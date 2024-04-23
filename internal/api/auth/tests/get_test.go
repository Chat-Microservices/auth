package tests

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/api/auth"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/service"
	serviceMocks "github.com/semho/chat-microservices/auth/internal/service/mocks"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestImplementation_Get(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService
	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		id      = gofakeit.Int64()
		name    = gofakeit.Name()
		email   = gofakeit.Email()
		role    = 1
		created = time.Now()
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		req     = &desc.GetRequest{
			Id: id,
		}

		res = &desc.UserResponse{
			Id: id,
			Detail: &desc.UserDetail{
				Name:  name,
				Email: email,
				Role:  desc.Role(role - 1).Enum(),
			},
			CreatedAt: timestamppb.New(created),
			UpdatedAt: timestamppb.New(created),
		}

		user = &model.User{
			ID: id,
			Detail: model.Detail{
				Name:  name,
				Email: email,
				Role:  role,
			},
			CreatedAt: created,
			UpdatedAt: created,
		}

		serviceError = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.UserResponse
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
				mock.GetMock.Expect(ctx, id).Return(user, nil)
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
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceError)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				authServiceMock := tt.authServiceMock(mc)
				api := authAPI.NewImplementation(authServiceMock)

				resHandler, err := api.Get(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			})
		})
	}
}
