package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, user *model.Detail, pass string) (int64, error) {
	id, err := s.authRepository.Create(ctx, user, pass)
	if err != nil {
		return 0, err
	}

	return id, nil
}
