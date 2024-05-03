package app

import (
	"context"
	accessAPI "github.com/semho/chat-microservices/auth/internal/api/access"
	authAPI "github.com/semho/chat-microservices/auth/internal/api/auth"
	loginAPI "github.com/semho/chat-microservices/auth/internal/api/login"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/client/db/pg"
	"github.com/semho/chat-microservices/auth/internal/client/db/transaction"
	"github.com/semho/chat-microservices/auth/internal/closer"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/config/env"
	"github.com/semho/chat-microservices/auth/internal/repository"
	accessRepository "github.com/semho/chat-microservices/auth/internal/repository/access"
	authRepository "github.com/semho/chat-microservices/auth/internal/repository/auth"
	loginRepository "github.com/semho/chat-microservices/auth/internal/repository/login"
	"github.com/semho/chat-microservices/auth/internal/service"
	accessService "github.com/semho/chat-microservices/auth/internal/service/access"
	authService "github.com/semho/chat-microservices/auth/internal/service/auth"
	loginService "github.com/semho/chat-microservices/auth/internal/service/login"
	"log"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient         db.Client
	txManger         db.TxManager
	authRepository   repository.AuthRepository
	loginRepository  repository.LoginRepository
	accessRepository repository.AccessRepository

	authService   service.AuthService
	loginService  service.LoginService
	accessService service.AccessService

	authImpl   *authAPI.Implementation
	loginImpl  *loginAPI.Implementation
	accessImpl *accessAPI.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %v", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) GetDBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to get db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) GetTxManager(ctx context.Context) db.TxManager {
	if s.txManger == nil {
		s.txManger = transaction.NewTransactionManager(s.GetDBClient(ctx).DB())
	}

	return s.txManger
}

func (s *serviceProvider) GetAuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.GetDBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) GetAuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.GetAuthRepository(ctx), s.GetTxManager(ctx))
	}

	return s.authService
}

func (s *serviceProvider) GetAuthImpl(ctx context.Context) *authAPI.Implementation {
	if s.authImpl == nil {
		s.authImpl = authAPI.NewImplementation(s.GetAuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) GetLoginRepository(ctx context.Context) repository.LoginRepository {
	if s.loginRepository == nil {
		s.loginRepository = loginRepository.NewRepository(s.GetDBClient(ctx))
	}

	return s.loginRepository
}

func (s *serviceProvider) GetLoginService(ctx context.Context) service.LoginService {
	if s.loginService == nil {
		s.loginService = loginService.NewService(s.GetLoginRepository(ctx), s.GetTxManager(ctx))
	}

	return s.loginService
}

func (s *serviceProvider) GetLoginImpl(ctx context.Context) *loginAPI.Implementation {
	if s.loginImpl == nil {
		s.loginImpl = loginAPI.NewImplementation(s.GetLoginService(ctx))
	}

	return s.loginImpl
}

func (s *serviceProvider) GetAccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepository.NewRepository(s.GetDBClient(ctx))
	}

	return s.accessRepository
}

func (s *serviceProvider) GetAccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.GetAccessRepository(ctx), s.GetTxManager(ctx))
	}

	return s.accessService
}

func (s *serviceProvider) GetAccessImpl(ctx context.Context) *accessAPI.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessAPI.NewImplementation(s.GetAccessService(ctx))
	}

	return s.accessImpl
}
