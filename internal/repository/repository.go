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
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type AccessRepository interface {
	AccessibleRoles(ctx context.Context) (map[string]int, error)
	GetListAccess(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Access, error)
	CreateAccess(ctx context.Context, roleId int, path string) (int64, error)
	DeleteAccess(ctx context.Context, id int64) error
	UpdateAccess(ctx context.Context, access *model.Access) error
}
