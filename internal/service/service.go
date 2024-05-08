package service

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

type AuthService interface {
	Get(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.Detail, pass string) (int64, error)
	Update(ctx context.Context, updateUser *model.UpdateUserRequest) error
	Delete(ctx context.Context, id int64) error
	GetListLogs(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Log, error)
}

type LoginService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpoint string) error
	AccessibleRoles(ctx context.Context) (map[string]int, error)
}
