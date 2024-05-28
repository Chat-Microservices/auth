package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/semho/chat-microservices/auth/internal/closer"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/interceptor"
	"github.com/semho/chat-microservices/auth/internal/logger"
	"github.com/semho/chat-microservices/auth/internal/metric"
	"github.com/semho/chat-microservices/auth/internal/tracing"
	descAccess "github.com/semho/chat-microservices/auth/pkg/access_v1"
	descAuth "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	descLogin "github.com/semho/chat-microservices/auth/pkg/login_v1"
	_ "github.com/semho/chat-microservices/auth/statik"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type App struct {
	servicesProvider *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatal(
				"failed to start grpc server: ", zap.Error(err),
			)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal(
				"failed to http swagger server: ", zap.Error(err),
			)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Fatal(
				"failed to start swagger server: ", zap.Error(err),
			)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusServer()
		if err != nil {
			logger.Fatal(
				"failed to start prometheus server: ", zap.Error(err),
			)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.InitLogger,
		a.InitMetrics,
		a.initServiceProvider,
		a.initTokenConfig,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

var configPath string
var logLevel string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.StringVar(&logLevel, "l", "info", "log level")
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) InitLogger(_ context.Context) error {
	flag.Parse()
	err := logger.InitDefault(logLevel)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) InitMetrics(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.servicesProvider = newServiceProvider()
	return nil
}

func (a *App) initTokenConfig(_ context.Context) error {
	a.servicesProvider.TokenConfig()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	fmt.Println("initGRPCServer")
	tracing.Init(logger.Logger(), "auth-service")

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		//grpc.Creds(credsService()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
				interceptor.MetricsInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)
	descAuth.RegisterAuthV1Server(a.grpcServer, a.servicesProvider.GetAuthImpl(ctx))
	descLogin.RegisterLoginV1Server(a.grpcServer, a.servicesProvider.GetLoginImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.servicesProvider.GetAccessImpl(ctx))

	return nil
}

func credsService() credentials.TransportCredentials {
	//Получаем текущую директорию
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	certFile := filepath.Join(dir, "tls", "service.pem")
	certKey := filepath.Join(dir, "tls", "service.key")

	creds, err := credentials.NewServerTLSFromFile(certFile, certKey)
	if err != nil {
		log.Fatalf("failed to load certificates: %v", err)
	}

	return creds
}

func (a *App) runGRPCServer() error {
	logger.Infof("Starting gRPC server on port: %s", a.servicesProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.servicesProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descAuth.RegisterAuthV1HandlerFromEndpoint(ctx, mux, a.servicesProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodHead,
				http.MethodOptions,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		},
	)

	a.httpServer = &http.Server{
		Addr:    a.servicesProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logger.Infof("Starting HTTP server on port: %s", a.servicesProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc(
		"/api.swagger.json",
		serveSwaggerFile("/api.swagger.json", a.servicesProvider.HTTPConfig().IpAddress()),
	)

	a.swaggerServer = &http.Server{
		Addr:    a.servicesProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Infof("Starting Swagger server on port: %s", a.servicesProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path, address string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		var data map[string]any
		err = json.Unmarshal(content, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["host"] = address

		modifiedContent, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(modifiedContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    a.servicesProvider.PrometheusConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	logger.Infof("Starting Prometheus server on port: %s", a.servicesProvider.PrometheusConfig().Address())

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
