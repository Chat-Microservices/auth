package accessService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s serv) GetListAccess(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Access, error) {
	var listAccess []*model.Access

	listAccess, err := s.accessRepository.GetListAccess(ctx, pageNumber, pageSize)

	if err != nil {
		return nil, err
	}

	return listAccess, nil
}
