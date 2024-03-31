package authService

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func (s *serv) GetListLogs(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Log, error) {
	var logs []*model.Log

	logs, err := s.authRepository.GetListLog(ctx, pageNumber, pageSize)

	if err != nil {
		return nil, err
	}

	return logs, nil
}
