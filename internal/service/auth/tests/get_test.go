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
	"time"
)

func Test_serv_Get(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		id      = gofakeit.Int64()
		name    = gofakeit.Name()
		email   = gofakeit.Email()
		role    = 1
		created = time.Now()

		res = &model.User{
			ID: id,
			Detail: model.Detail{
				Name:  name,
				Email: email,
				Role:  role,
			},
			CreatedAt: created,
			UpdatedAt: created,
		}

		qName = "CreateLog"
		qRow  = "GET SELECT..."
		query = db.Query{
			Name:     qName,
			QueryRow: qRow,
		}

		repoErr = fmt.Errorf("repo error")
		log     = &model.Log{
			Action:   qName,
			EntityId: id,
			Query:    qRow,
		}
		repoErrLog = fmt.Errorf("repo error log")
	)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		authRepositoryMock authRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: res,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, query, nil)
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
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, query, repoErr)
				return mock
			},
		},
		{
			name: "create log error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErrLog,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, query, nil)
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

			resHandler, err := service.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
