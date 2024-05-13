package accessService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s serv) UpdateAccess(ctx context.Context, access *model.Access) error {
	err := s.accessRepository.UpdateAccess(ctx, access)
	if err != nil {
		return err
	}

	return nil
}
