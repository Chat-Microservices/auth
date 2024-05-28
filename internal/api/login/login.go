package loginAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"go.uber.org/zap"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	logger.Info("username", zap.String("username", req.GetUsername()))
	refreshToken, err := i.loginService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	logger.Info("refreshToken", zap.String("refreshToken", refreshToken))

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
