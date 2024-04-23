package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/api/auth"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/service"
	serviceMocks "github.com/semho/chat-microservices/auth/internal/service/mocks"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestImplementation_Update(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		req = &desc.UpdateRequest{
			Id: id,
			Info: &desc.UpdateUserInfo{
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
			},
		}
		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.UpdateMock.Expect(ctx, converter.ToAuthUpdateUserFromDesc(req)).Return(nil)
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
				mock.UpdateMock.Expect(ctx, converter.ToAuthUpdateUserFromDesc(req)).Return(serviceErr)
				return mock
			},
		},
		{
			name: "error without id",
			args: args{
				ctx: ctx,
				req: &desc.UpdateRequest{
					Id: 0,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Id must be provided"),
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
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

				resHandler, err := api.Update(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			})
		})
	}
}
