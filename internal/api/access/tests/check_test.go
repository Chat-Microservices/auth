package tests

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/api/access"
	"github.com/semho/chat-microservices/auth/internal/service"
	serviceMocks "github.com/semho/chat-microservices/auth/internal/service/mocks"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestImplementation_Check(t *testing.T) {
	t.Parallel()
	type accessServiceMockFunc func(mc *minimock.Controller) service.AccessService
	type args struct {
		ctx context.Context
		req *desc.CheckRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.CheckRequest{
			EndpointAddress: "endpointAddress",
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
				mock.CheckMock.Expect(ctx, req.EndpointAddress).Return(nil)
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
				mock.CheckMock.Expect(ctx, req.EndpointAddress).Return(serviceError)
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
						accessServiceMock := tt.accessServiceMock(mc)
						api := accessAPI.NewImplementation(accessServiceMock)

						resHandler, err := api.Check(tt.args.ctx, tt.args.req)
						require.Equal(t, tt.err, err)
						require.Equal(t, tt.want, resHandler)
					},
				)
			},
		)
	}
}
