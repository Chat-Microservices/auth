package loginAPI

import (
	"context"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"log"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	log.Printf("username: %s", req.GetUsername())
	refreshToken, err := i.loginService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	log.Printf("get refreshToken: %s", refreshToken)

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
