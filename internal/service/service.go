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
}
