package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/api/auth"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/service"
	serviceMocks "github.com/semho/chat-microservices/auth/internal/service/mocks"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestImplementation_Create(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = 1

		password = gofakeit.Color()

		serviceErr = fmt.Errorf("service error")
		detail     = &desc.UserDetail{
			Name:  name,
			Email: email,
			Role:  desc.Role(role - 1).Enum(),
		}
		correctPass = &desc.UserPassword{
			Password:        password,
			PasswordConfirm: password,
		}
		inCorrectPass = &desc.UserPassword{
			Password:        password,
			PasswordConfirm: "incorrect_confirmation",
		}
		req = &desc.CreateRequest{
			Detail:   detail,
			Password: correctPass,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(
					ctx, &model.Detail{
						Name:  name,
						Email: email,
						Role:  1,
					}, password,
				).Return(id, nil)
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
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(
					ctx, &model.Detail{
						Name:  name,
						Email: email,
						Role:  1,
					}, password,
				).Return(0, serviceErr)
				return mock
			},
		},
		{
			name: "error with nil details",
			args: args{
				ctx: ctx,
				req: &desc.CreateRequest{
					Detail:   nil, // Detail is nil
					Password: correctPass,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Detail and Password must be provided"),
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				return mock
			},
		},
		{
			name: "error with nil password",
			args: args{
				ctx: ctx,
				req: &desc.CreateRequest{
					Detail:   detail,
					Password: nil, // Password is nil
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Detail and Password must be provided"),
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				return mock
			},
		},
		{
			name: "error with mismatched password",
			args: args{
				ctx: ctx,
				req: &desc.CreateRequest{
					Detail:   detail,
					Password: inCorrectPass,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Password and Password Confirm do not match"),
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				return mock
			},
		},
	}
	initLogger()
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				authServiceMock := tt.authServiceMock(mc)
				api := authAPI.NewImplementation(authServiceMock)

				resHandler, err := api.Create(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
