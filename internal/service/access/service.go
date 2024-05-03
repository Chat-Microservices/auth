package accessService

import (
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/repository"
	"github.com/semho/chat-microservices/auth/internal/service"
)

type serv struct {
	accessRepository repository.AccessRepository
	txManager        db.TxManager
}

func NewService(accessRepository repository.AccessRepository, txManager db.TxManager) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		txManager:        txManager,
	}
}
