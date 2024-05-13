package accessService

import (
	"context"
)

func (s serv) DeleteAccess(ctx context.Context, id int64) error {
	err := s.accessRepository.DeleteAccess(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
