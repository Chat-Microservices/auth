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

func TestImplementation_Login(t *testing.T) {
	t.Parallel()
	type loginServiceMockFunc func(mc *minimock.Controller) service.LoginService
	type args struct {
		ctx context.Context
		req *desc.LoginRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.LoginRequest{
			Username: "username",
			Password: "password",
		}
		res = &desc.LoginResponse{
			RefreshToken: "refreshToken",
		}
		serviceError = errors.New("service error")
	)

	tests := []struct {
		name             string
		args             args
		want             *desc.LoginResponse
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
				mock.LoginMock.Expect(ctx, req.Username, req.Password).Return(res.RefreshToken, nil)
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
				mock.LoginMock.Expect(ctx, req.Username, req.Password).Return("", serviceError)
				return mock
			},
		},
	}
	initLogger()
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt := tt
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						loginServiceMock := tt.loginServiceMock(mc)
						api := loginAPI.NewImplementation(loginServiceMock)

						resHandler, err := api.Login(tt.args.ctx, tt.args.req)
						require.Equal(t, tt.err, err)
						require.Equal(t, tt.want, resHandler)
					},
				)
			},
		)
	}
}
