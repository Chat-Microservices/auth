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
	authService "github.com/semho/chat-microservices/auth/internal/service/auth"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImplementation_Create(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository

	type args struct {
		ctx      context.Context
		req      *model.Detail
		password string
	}

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = 1

		password = gofakeit.Color()

		repoErr = fmt.Errorf("repo error")
		req     = &model.Detail{
			Name:  name,
			Email: email,
			Role:  role,
		}

		qName = "CreateLog"
		qRow  = "INSERT INTO..."
		query = db.Query{
			Name:     qName,
			QueryRow: qRow,
		}

		log = &model.Log{
			Action:   qName,
			EntityId: id,
			Query:    qRow,
		}
		repoErrLog = fmt.Errorf("repo error log")
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		authRepositoryMock authRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:      ctx,
				req:      req,
				password: password,
			},
			want: id,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				//mock.CreateMock.Expect(ctx, req, password).Return(id, query, nil)
				//expectedLogEntry := converter.ToAuthLogFromQuery(query, id)
				//mock.CreateLogMock.Expect(ctx, expectedLogEntry).Return(nil)
				mock.CreateMock.Set(
					func(ctx context.Context, req *model.Detail, password string) (int64, db.Query, error) {
						if ctx == nil {
							return 0, db.Query{}, fmt.Errorf("context is nil")
						}
						if req == nil {
							return 0, db.Query{}, fmt.Errorf("req is nil")
						}
						if password == "" {
							return 0, db.Query{}, fmt.Errorf("password is empty")
						}
						return id, query, nil
					},
				)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx:      ctx,
				req:      req,
				password: password,
			},
			want: 0,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				//mock.CreateMock.Expect(ctx, req, password).Return(0, query, repoErr)
				mock.CreateMock.Set(
					func(ctx context.Context, req *model.Detail, password string) (int64, db.Query, error) {
						return 0, db.Query{}, repoErr
					},
				)
				return mock
			},
		},
		{
			name: "create log error case",
			args: args{
				ctx:      ctx,
				req:      req,
				password: password,
			},
			want: 0,
			err:  repoErrLog,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				//mock.CreateMock.Expect(ctx, req, password).Return(id, query, nil)
				mock.CreateMock.Set(
					func(ctx context.Context, req *model.Detail, password string) (int64, db.Query, error) {
						return id, query, nil
					},
				)
				mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				fmt.Println("[TEST DEBUG] Starting test case:", tt.name)
				authRepoMock := tt.authRepositoryMock(mc)
				fmt.Println("[TEST DEBUG] authRepoMock:", authRepoMock)

				service := authService.NewService(authRepoMock, txManagerMock())
				fmt.Println("[TEST DEBUG] service instance:", service)

				if tt.args.req == nil {
					t.Fatal("req is nil")
				}
				resHandler, err := service.Create(tt.args.ctx, tt.args.req, tt.args.password)
				fmt.Println("[TEST DEBUG] resHandler:", resHandler, "error:", err)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
