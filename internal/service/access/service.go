package accessService

import (
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/repository"
	"github.com/semho/chat-microservices/auth/internal/service"
)

type serv struct {
	accessRepository repository.AccessRepository
	txManager        db.TxManager
	accessibleRoles  map[string]int
	tokenConfig      config.TokenConfig
}

func NewService(
	accessRepository repository.AccessRepository,
	txManager db.TxManager,
	tokenConfig config.TokenConfig,
) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		txManager:        txManager,
		tokenConfig:      tokenConfig,
	}
}
