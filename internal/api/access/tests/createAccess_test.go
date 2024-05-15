package tests

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/api/access"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/service"
	serviceMocks "github.com/semho/chat-microservices/auth/internal/service/mocks"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestImplementation_CreateAccess(t *testing.T) {
	t.Parallel()
	type accessServiceMockFunc func(mc *minimock.Controller) service.AccessService
	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.CreateRequest{
			RoleId: 1,
			Path:   "test",
		}

		res = &desc.CreateResponse{
			Id: 1,
		}

		serviceError = errors.New("service error")
	)
	tests := []struct {
		name              string
		args              args
		want              *desc.CreateResponse
		err               error
		accessServiceMock accessServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.CreateAccessMock.Expect(
					ctx, &model.Access{
						RoleId: int(req.RoleId + 1),
						Path:   req.Path,
					},
				).Return(res.Id, nil)
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
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.CreateAccessMock.Expect(
					ctx, &model.Access{
						RoleId: int(req.RoleId + 1),
						Path:   req.Path,
					},
				).Return(0, serviceError)
				return mock
			},
		},
		{
			name: "error with empty path",
			args: args{
				ctx: ctx,
				req: &desc.CreateRequest{
					RoleId: 1,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Role and Path must be provided"),
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				return serviceMocks.NewAccessServiceMock(mc)
			},
		},
		{
			name: "error with zero role",
			args: args{
				ctx: ctx,
				req: &desc.CreateRequest{
					Path: "test",
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Role and Path must be provided"),
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				return serviceMocks.NewAccessServiceMock(mc)
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
						accessServiceMock := tt.accessServiceMock(mc)
						api := accessAPI.NewImplementation(accessServiceMock)

						resHandler, err := api.CreateAccess(tt.args.ctx, tt.args.req)
						require.Equal(t, tt.err, err)
						require.Equal(t, tt.want, resHandler)
					},
				)
			},
		)
	}
}
