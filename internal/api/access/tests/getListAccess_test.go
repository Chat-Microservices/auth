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
	"testing"
)

func TestImplementation_GetListAccess(t *testing.T) {
	t.Parallel()
	type accessServiceMockFunc func(mc *minimock.Controller) service.AccessService
	type args struct {
		ctx context.Context
		req *desc.GetListAccessRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.GetListAccessRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		res = &desc.ListAccessResponse{
			List: []*desc.Access{
				{
					Id:     1,
					RoleId: 1,
					Path:   "path",
				},
			},
		}

		expectedAccess = []*model.Access{
			{
				ID:     1,
				RoleId: 1 + 1,
				Path:   "path",
			},
		}

		serviceError = errors.New("service error")
	)
	tests := []struct {
		name              string
		args              args
		want              *desc.ListAccessResponse
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
				mock.GetListAccessMock.Expect(ctx, req.GetPageNumber(), req.GetPageSize()).Return(expectedAccess, nil)
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
				mock.GetListAccessMock.Expect(ctx, req.GetPageNumber(), req.GetPageSize()).Return(nil, serviceError)
				return mock
			},
		},
		{
			name: "page number 0",
			args: args{
				ctx: ctx,
				req: &desc.GetListAccessRequest{
					PageNumber: 0,
					PageSize:   1,
				},
			},
			want: res,
			err:  nil,
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.GetListAccessMock.Expect(ctx, 1, req.GetPageSize()).Return(expectedAccess, nil)
				return mock
			},
		},
		{
			name: "page size 0",
			args: args{
				ctx: ctx,
				req: &desc.GetListAccessRequest{
					PageNumber: 1,
					PageSize:   0,
				},
			},
			want: res,
			err:  nil,
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.GetListAccessMock.Expect(ctx, req.GetPageNumber(), 5).Return(expectedAccess, nil)
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
						accessServiceMock := tt.accessServiceMock(mc)
						api := accessAPI.NewImplementation(accessServiceMock)

						resHandler, err := api.GetListAccess(tt.args.ctx, tt.args.req)
						require.Equal(t, tt.err, err)
						require.Equal(t, tt.want, resHandler)
					},
				)
			},
		)
	}
}
