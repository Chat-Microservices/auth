package accessAPI

import (
	"github.com/semho/chat-microservices/auth/internal/service"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
