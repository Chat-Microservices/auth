package loginAPI

import (
	"context"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"log"
)

func (i *Implementation) GetAccessToken(
	ctx context.Context,
	req *desc.GetAccessTokenRequest,
) (*desc.GetAccessTokenResponse, error) {
	log.Printf("refresh token: %s", req.GetRefreshToken())
	accessToken, err := i.loginService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	log.Printf("get accessToken: %s", accessToken)

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
