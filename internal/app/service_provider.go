package app

import (
	"context"
	authAPI "github.com/semho/chat-microservices/auth/internal/api/auth"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/client/db/pg"
	"github.com/semho/chat-microservices/auth/internal/client/db/transaction"
	"github.com/semho/chat-microservices/auth/internal/closer"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/config/env"
	"github.com/semho/chat-microservices/auth/internal/repository"
	authRepository "github.com/semho/chat-microservices/auth/internal/repository/auth"
	"github.com/semho/chat-microservices/auth/internal/service"
	authService "github.com/semho/chat-microservices/auth/internal/service/auth"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManger       db.TxManager
	authRepository repository.AuthRepository

	authService service.AuthService

	authImpl *authAPI.Implementation
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
