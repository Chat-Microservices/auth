package repository

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

type AuthRepository interface {
	Get(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, detail *model.Detail, pass string) (int64, error)
	Update(ctx context.Context, updateUser *model.UpdateUserRequest) error
	Delete(ctx context.Context, id int64) error
}
