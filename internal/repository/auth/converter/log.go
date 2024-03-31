package converter

import (
	"github.com/semho/chat-microservices/auth/internal/model"
	modelRepo "github.com/semho/chat-microservices/auth/internal/repository/auth/model"
)

func ToLogFromRepo(logs []modelRepo.Log) []*model.Log {
	convertedLogs := make([]*model.Log, len(logs))
	for i, log := range logs {
		convertedLogs[i] = &model.Log{
			ID:        log.ID,
			Action:    log.Action,
			EntityId:  log.EntityId,
			Query:     log.Query,
			CreatedAt: log.CreatedAt,
			UpdatedAt: log.UpdatedAt,
		}
	}
	return convertedLogs
}
