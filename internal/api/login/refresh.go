package loginAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"go.uber.org/zap"
)

func (i *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	logger.Info("oldRefreshToken", zap.String("oldRefreshToken", req.GetOldRefreshToken()))
	refreshToken, err := i.loginService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	logger.Info("get refreshToken", zap.String("refreshToken", refreshToken))

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
