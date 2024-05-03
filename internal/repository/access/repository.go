package loginRepository

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/repository"
)

const (
	tableName = "users"
)

type repo struct {
	db db.Client
}

func (r repo) Check(ctx context.Context, endpoint string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}
