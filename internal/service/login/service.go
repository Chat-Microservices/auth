package loginService

import (
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/repository"
	"github.com/semho/chat-microservices/auth/internal/service"
)

type serv struct {
	loginRepository repository.LoginRepository
	txManager       db.TxManager
	tokenConfig     config.TokenConfig
}

func NewService(
	loginRepository repository.LoginRepository,
	txManager db.TxManager,
	tokenConfig config.TokenConfig,
) service.LoginService {
	return &serv{
		loginRepository: loginRepository,
		txManager:       txManager,
		tokenConfig:     tokenConfig,
	}
}
