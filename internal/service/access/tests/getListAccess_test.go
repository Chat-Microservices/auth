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
	"github.com/semho/chat-microservices/auth/internal/service/access"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_serv_GetListAccess(t *testing.T) {
	t.Parallel()
	type accessRepositoryMockFunc func(mc *minimock.Controller) repository.AccessRepository
	type args struct {
		ctx        context.Context
		pageNumber uint64
		pageSize   uint64
	}

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	tokenConfigMock := func() config.TokenConfig {
		return &MockTokenConfig{}
	}

	var (
		ctx        = context.Background()
		mc         = minimock.NewController(t)
		pageNumber = gofakeit.Uint64()
		pageSize   = gofakeit.Uint64()
		res        = []*model.Access{
			{
				ID:     gofakeit.Int64(),
				RoleId: 1,
				Path:   gofakeit.URL(),
			},
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name                 string
		args                 args
		want                 []*model.Access
		err                  error
		accessRepositoryMock accessRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:        ctx,
				pageNumber: pageNumber,
				pageSize:   pageSize,
			},
			want: res,
			err:  nil,
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repoMocks.NewAccessRepositoryMock(mc)
				mock.GetListAccessMock.Expect(ctx, pageNumber, pageSize).Return(res, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx:        ctx,
				pageNumber: pageNumber,
				pageSize:   pageSize,
			},
			want: nil,
			err:  repoErr,
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repoMocks.NewAccessRepositoryMock(mc)
				mock.GetListAccessMock.Expect(ctx, pageNumber, pageSize).Return(nil, repoErr)
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

				resHandler, err := service.GetListAccess(tt.args.ctx, tt.args.pageNumber, tt.args.pageSize)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
