package loginAPI

import (
	"github.com/semho/chat-microservices/auth/internal/service"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
)

type Implementation struct {
	desc.UnimplementedLoginV1Server
	loginService service.LoginService
}

func NewImplementation(loginService service.LoginService) *Implementation {
	return &Implementation{
		loginService: loginService,
	}
}
