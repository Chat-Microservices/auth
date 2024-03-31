package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/converter"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	var query db.Query
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		query, errTx = s.authRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.authRepository.CreateLog(ctx, converter.ToAuthLogFromQuery(query, id))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
