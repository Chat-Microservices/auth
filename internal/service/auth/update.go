package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, updateUser *model.UpdateUserRequest) error {
	err := s.authRepository.Update(ctx, updateUser)
	if err != nil {
		return err
	}

	return nil
}
