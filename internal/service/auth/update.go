package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, updateUser *model.UpdateUserRequest) error {
	var query db.Query
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		query, errTx = s.authRepository.Update(ctx, updateUser)
		if errTx != nil {
			return errTx
		}

		errTx = s.authRepository.CreateLog(ctx, converter.ToAuthLogFromQuery(query, updateUser.ID))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
