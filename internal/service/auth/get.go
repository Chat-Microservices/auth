package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User
	var query db.Query
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, query, errTx = s.authRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.authRepository.CreateLog(ctx, converter.ToAuthLogFromQuery(query, id))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
