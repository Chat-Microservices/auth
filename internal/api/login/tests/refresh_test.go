package tests

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/api/login"
	"github.com/semho/chat-microservices/auth/internal/service"
	serviceMocks "github.com/semho/chat-microservices/auth/internal/service/mocks"
	desc "github.com/semho/chat-microservices/auth/pkg/login_v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImplementation_GetRefreshToken(t *testing.T) {
	t.Parallel()
	type loginServiceMockFunc func(mc *minimock.Controller) service.LoginService
	type args struct {
		ctx context.Context
		req *desc.GetRefreshTokenRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.GetRefreshTokenRequest{
			OldRefreshToken: "oldRefreshToken",
		}
		res = &desc.GetRefreshTokenResponse{
			RefreshToken: "refreshToken",
		}
		serviceError = errors.New("service error")
	)
	tests := []struct {
		name             string
		args             args
		want             *desc.GetRefreshTokenResponse
		err              error
		loginServiceMock loginServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			loginServiceMock: func(mc *minimock.Controller) service.LoginService {
				mock := serviceMocks.NewLoginServiceMock(mc)
				mock.GetRefreshTokenMock.Expect(ctx, req.OldRefreshToken).Return(res.RefreshToken, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceError,
			loginServiceMock: func(mc *minimock.Controller) service.LoginService {
				mock := serviceMocks.NewLoginServiceMock(mc)
				mock.GetRefreshTokenMock.Expect(ctx, req.OldRefreshToken).Return("", serviceError)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt := tt
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						loginServiceMock := tt.loginServiceMock(mc)
						api := loginAPI.NewImplementation(loginServiceMock)

						resHandler, err := api.GetRefreshToken(tt.args.ctx, tt.args.req)
						require.Equal(t, tt.err, err)
						require.Equal(t, tt.want, resHandler)
					},
				)
			},
		)
	}
}
