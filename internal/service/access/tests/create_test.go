package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/repository"
	repoMocks "github.com/semho/chat-microservices/auth/internal/repository/mocks"
	accessService "github.com/semho/chat-microservices/auth/internal/service/access"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_serv_CreateAccess(t *testing.T) {
	t.Parallel()
	type accessRepositoryMockFunc func(mc *minimock.Controller) repository.AccessRepository
	type args struct {
		ctx context.Context
		req *model.Access
	}

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	tokenConfigMock := func() config.TokenConfig {
		return &MockTokenConfig{}
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = gofakeit.Int64()
		req = &model.Access{
			RoleId: 1,
			Path:   "path",
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name                 string
		args                 args
		want                 int64
		err                  error
		accessRepositoryMock accessRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repoMocks.NewAccessRepositoryMock(mc)
				mock.CreateAccessMock.Expect(ctx, req.RoleId, req.Path).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repoMocks.NewAccessRepositoryMock(mc)
				mock.CreateAccessMock.Expect(ctx, req.RoleId, req.Path).Return(0, repoErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				accessRepoMock := tt.accessRepositoryMock(mc)
				service := accessService.NewService(accessRepoMock, txManagerMock(), tokenConfigMock())

				if tt.args.req == nil {
					t.Fatal("req is nil")
				}
				resHandler, err := service.CreateAccess(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
