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
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestImplementation_UpdateAccess(t *testing.T) {
	t.Parallel()
	type accessServiceMockFunc func(mc *minimock.Controller) service.AccessService
	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.UpdateRequest{
			Id:     1,
			RoleId: wrapperspb.Int32(1),
			Path:   wrapperspb.String("path"),
		}

		serviceError = errors.New("service error")
	)
	tests := []struct {
		name              string
		args              args
		want              *emptypb.Empty
		err               error
		accessServiceMock accessServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.UpdateAccessMock.Expect(
					ctx, &model.Access{
						ID:     req.GetId(),
						RoleId: int(req.GetRoleId().GetValue()) + 1,
						Path:   req.GetPath().GetValue(),
					},
				).Return(nil)
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
				mock.UpdateAccessMock.Expect(
					ctx, &model.Access{
						ID:     req.GetId(),
						RoleId: int(req.GetRoleId().GetValue()) + 1,
						Path:   req.GetPath().GetValue(),
					},
				).Return(serviceError)
				return mock
			},
		},
		{
			name: "error with empty id",
			args: args{
				ctx: ctx,
				req: &desc.UpdateRequest{},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Id must be provided"),
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				return serviceMocks.NewAccessServiceMock(mc)
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
						accessServiceMock := tt.accessServiceMock(mc)
						api := accessAPI.NewImplementation(accessServiceMock)

						resHandler, err := api.UpdateAccess(tt.args.ctx, tt.args.req)
						require.Equal(t, tt.err, err)
						require.Equal(t, tt.want, resHandler)
					},
				)
			},
		)
	}
}
