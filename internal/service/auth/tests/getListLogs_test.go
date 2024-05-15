package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/repository"
	repoMocks "github.com/semho/chat-microservices/auth/internal/repository/mocks"
	"github.com/semho/chat-microservices/auth/internal/service/auth"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_serv_GetListLogs(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository

	type args struct {
		ctx        context.Context
		pageNumber int64
		pageSize   int64
	}

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}
	var (
		ctx        = context.Background()
		mc         = minimock.NewController(t)
		pageNumber = gofakeit.Uint64()
		pageSize   = gofakeit.Uint64()
		created    = time.Now()

		res = []*model.Log{
			{
				ID:        gofakeit.Int64(),
				Action:    gofakeit.BeerName(),
				EntityId:  gofakeit.Int64(),
				Query:     gofakeit.City(),
				CreatedAt: created,
				UpdatedAt: created,
			},
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name               string
		args               args
		want               []*model.Log
		err                error
		authRepositoryMock authRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:        ctx,
				pageNumber: int64(pageNumber),
				pageSize:   int64(pageSize),
			},
			want: res,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetListLogMock.Expect(ctx, pageNumber, pageSize).Return(res, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx:        ctx,
				pageNumber: int64(pageNumber),
				pageSize:   int64(pageSize),
			},
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetListLogMock.Expect(ctx, pageNumber, pageSize).Return(nil, repoErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				authRepoMock := tt.authRepositoryMock(mc)
				service := authService.NewService(authRepoMock, txManagerMock())

				resHandler, err := service.GetListLogs(
					tt.args.ctx,
					uint64(tt.args.pageNumber),
					uint64(tt.args.pageSize),
				)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
