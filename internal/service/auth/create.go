package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, user *model.Detail, pass string) (int64, error) {
	var id int64
	var query db.Query
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, query, errTx = s.authRepository.Create(ctx, user, pass)
		if errTx != nil {
			return errTx
		}

		errTx = s.authRepository.CreateLog(ctx, converter.ToAuthLogFromQuery(query, id))

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
