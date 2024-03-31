package authAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.UserResponse, error) {
	log.Printf("User id: %d", req.GetId())
	userObj, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("get user: %v", userObj)

	return converter.ToAuthFromService(userObj), nil
}
