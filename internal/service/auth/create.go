package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, user *model.Detail, pass string) (int64, error) {
	//id, err := s.authRepository.Create(ctx, user, pass)
	//if err != nil {
	//	return 0, err
	//}
	//
	//return id, nil

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.authRepository.Create(ctx, user, pass)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.authRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
