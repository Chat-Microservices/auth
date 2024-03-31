package converter

import (
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToAuthLogFromQuery(q db.Query, id int64) *model.Log {
	return &model.Log{
		Action:   q.Name,
		EntityId: id,
		Query:    q.QueryRow,
	}
}

func ToLogFromService(logs []*model.Log) *desc.LogsResponse {
	logsResponses := make([]*desc.Log, len(logs))
	for i, log := range logs {
		logsResponses[i] = &desc.Log{
			Id:        log.ID,
			Action:    log.Action,
			EntityId:  log.EntityId,
			Query:     log.Query,
			CreatedAt: timestamppb.New(log.CreatedAt),
			UpdatedAt: timestamppb.New(log.UpdatedAt),
		}

	}
	return &desc.LogsResponse{
		Logs: logsResponses,
	}
}
