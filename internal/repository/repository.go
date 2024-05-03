package repository

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
)

type AuthRepository interface {
	Get(ctx context.Context, id int64) (*model.User, db.Query, error)
	Create(ctx context.Context, detail *model.Detail, pass string) (int64, db.Query, error)
	Update(ctx context.Context, updateUser *model.UpdateUserRequest) (db.Query, error)
	Delete(ctx context.Context, id int64) (db.Query, error)
	CreateLog(ctx context.Context, logger *model.Log) error
	GetListLog(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Log, error)
}

type LoginRepository interface {
	Login(ctx context.Context, username string, password string) (*model.User, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessRepository interface {
	Check(ctx context.Context, endpoint string) error
}
