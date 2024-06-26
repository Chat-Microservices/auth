package authAPI

import (
	"github.com/semho/chat-microservices/auth/internal/service"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
