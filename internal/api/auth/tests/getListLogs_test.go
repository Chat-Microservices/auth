package tests

import (
	"context"
	"errors"
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

func TestImplementation_GetListLogs(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService
	type args struct {
		ctx context.Context
		req *desc.GetListLogsRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.GetListLogsRequest{
			PageNumber: 1,
			PageSize:   1,
		}

		created = time.Now()

		res = &desc.LogsResponse{
			Logs: []*desc.Log{
				{
					Id:        1,
					Action:    "test",
					EntityId:  1,
					Query:     "test",
					CreatedAt: timestamppb.New(created),
					UpdatedAt: timestamppb.New(created),
				},
			},
		}
		expectedLogs = []*model.Log{
			{
				ID:        1,
				Action:    "test",
				EntityId:  1,
				Query:     "test",
				CreatedAt: created,
				UpdatedAt: created,
			},
		}
		serviceError = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.LogsResponse
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
				mock.GetListLogsMock.Expect(ctx, req.GetPageNumber(), req.GetPageSize()).Return(expectedLogs, nil)
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
				mock.GetListLogsMock.Expect(ctx, req.GetPageNumber(), req.GetPageSize()).Return(nil, serviceError)
				return mock
			},
		},
		{
			name: "page number 0",
			args: args{
				ctx: ctx,
				req: &desc.GetListLogsRequest{
					PageNumber: 0,
					PageSize:   1,
				},
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.GetListLogsMock.Expect(ctx, 1, req.GetPageSize()).Return(expectedLogs, nil)
				return mock
			},
		},
		{
			name: "page size 0",
			args: args{
				ctx: ctx,
				req: &desc.GetListLogsRequest{
					PageNumber: 1,
					PageSize:   0,
				},
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.GetListLogsMock.Expect(ctx, req.GetPageNumber(), 5).Return(expectedLogs, nil)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			authServiceMock := tt.authServiceMock(mc)
			api := authAPI.NewImplementation(authServiceMock)

			resHandler, err := api.GetListLogs(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
