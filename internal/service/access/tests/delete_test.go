package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/repository"
	repoMocks "github.com/semho/chat-microservices/auth/internal/repository/mocks"
	"github.com/semho/chat-microservices/auth/internal/service/access"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_serv_DeleteAccess(t *testing.T) {
	t.Parallel()
	type accessRepositoryMockFunc func(mc *minimock.Controller) repository.AccessRepository
	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
	)

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	tokenConfigMock := func() config.TokenConfig {
		return &MockTokenConfig{}
	}

	tests := []struct {
		name                 string
		args                 args
		err                  error
		accessRepositoryMock accessRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repoMocks.NewAccessRepositoryMock(mc)
				mock.DeleteAccessMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErr,
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repoMocks.NewAccessRepositoryMock(mc)
				mock.DeleteAccessMock.Expect(ctx, id).Return(repoErr)
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

				err := service.DeleteAccess(tt.args.ctx, tt.args.id)
				require.Equal(t, tt.err, err)
			},
		)
	}
}
