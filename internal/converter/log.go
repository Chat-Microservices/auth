package converter

import (
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
)

func ToAuthLogFromQuery(q db.Query, id int64) *model.Log {
	return &model.Log{
		Action:   q.Name,
		EntityId: id,
		Query:    q.QueryRow,
	}
}
