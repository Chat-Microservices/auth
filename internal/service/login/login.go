package loginService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s serv) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", status.Error(codes.InvalidArgument, "Invalid request: username and password must be provided")
	}

	user, err := s.loginRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(user.Password, password) {
		return "", status.Error(codes.Unauthenticated, "Invalid username or password")
	}

	refreshTokenSecretKey, refreshTokenExpiration := s.tokenConfig.RefreshData()
	
	refreshToken, err := utils.GenerateToken(
		user.Detail,
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", status.Errorf(codes.Aborted, "Failed to generate refresh token")
	}

	return refreshToken, nil
}
