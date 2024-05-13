package accessService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) CreateAccess(ctx context.Context, access *model.Access) (int64, error) {
	var id int64

	id, err := s.accessRepository.CreateAccess(ctx, access.RoleId, access.Path)
	if err != nil {
		return 0, err
	}

	return id, nil
}
