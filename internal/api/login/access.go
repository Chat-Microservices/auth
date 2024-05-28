package loginAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"go.uber.org/zap"
)

func (i *Implementation) GetAccessToken(
	ctx context.Context,
	req *desc.GetAccessTokenRequest,
) (*desc.GetAccessTokenResponse, error) {
	logger.Info("refreshToken", zap.String("refreshToken", req.GetRefreshToken()))
	accessToken, err := i.loginService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	logger.Info("get accessToken", zap.String("accessToken", accessToken))

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
