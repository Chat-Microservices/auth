package loginService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", status.Error(codes.InvalidArgument, "Invalid request: username and password must be provided")
	}

	refreshTokenSecretKey, _ := s.tokenConfig.RefreshData()
	accessTokenSecretKey, accessTokenExpiration := s.tokenConfig.AccessData()

	claims, err := utils.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}
	accessToken, err := utils.GenerateToken(
		model.Detail{
			Name:  claims.Username,
			Role:  claims.Role,
			Email: claims.Email,
		}, []byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
