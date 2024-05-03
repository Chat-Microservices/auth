package loginRepository

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/repository"
)

const (
	tableName = "users"
)

type repo struct {
	db db.Client
}

func (r repo) Login(ctx context.Context, username string, password string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db db.Client) repository.LoginRepository {
	return &repo{db: db}
}
