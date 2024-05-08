package loginService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	if oldRefreshToken == "" {
		return "", status.Error(codes.InvalidArgument, "Invalid request: username and password must be provided")
	}

	refreshTokenSecretKey, refreshTokenExpiration := s.tokenConfig.RefreshData()

	claims, err := utils.VerifyToken(oldRefreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(
		model.Detail{
			Name:  claims.Username,
			Role:  claims.Role,
			Email: claims.Email,
		}, []byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
