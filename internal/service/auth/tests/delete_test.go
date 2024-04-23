package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/repository"
	repoMocks "github.com/semho/chat-microservices/auth/internal/repository/mocks"
	authService "github.com/semho/chat-microservices/auth/internal/service/auth"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_serv_Delete(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		qName = "CreateLog"
		qRow  = "UPDATE..."
		query = db.Query{
			Name:     qName,
			QueryRow: qRow,
		}
		log = &model.Log{
			Action:   qName,
			EntityId: id,
			Query:    qRow,
		}

		repoErr    = fmt.Errorf("repo error")
		repoErrLog = fmt.Errorf("repo error log")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		authRepositoryMock authRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(query, nil)
				expectedLogEntry := converter.ToAuthLogFromQuery(query, id)
				mock.CreateLogMock.Expect(ctx, expectedLogEntry).Return(nil)
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
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(query, repoErr)
				return mock
			},
		},
		{
			name: "create log error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErrLog,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			authRepoMock := tt.authRepositoryMock(mc)
			service := authService.NewService(authRepoMock, txManagerMock())

			err := service.Delete(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
		})
	}
}
