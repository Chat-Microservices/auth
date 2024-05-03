package loginAPI

import (
	"context"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"log"
)

func (i *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	log.Printf("refresh token: %s", req.GetRefreshToken())
	refreshToken, err := i.loginService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	log.Printf("get refreshToken: %s", refreshToken)

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
