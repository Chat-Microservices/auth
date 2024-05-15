package converter

import (
	"github.com/semho/chat-microservices/auth/internal/model"
	modelRepo "github.com/semho/chat-microservices/auth/internal/repository/access/model"
)

func ToMapAccessFromRepo(access []modelRepo.Access) map[string]int {
	accessibleRoles := make(map[string]int)

	for _, onRow := range access {
		accessibleRoles[onRow.Path] = int(onRow.RoleId)
	}

	return accessibleRoles
}

func ToAccessListFromRepo(listAccess []modelRepo.Access) []*model.Access {
	convertedAccess := make([]*model.Access, len(listAccess))
	for i, log := range listAccess {
		convertedAccess[i] = &model.Access{
			ID:     log.ID,
			RoleId: int(log.RoleId),
			Path:   log.Path,
		}
	}
	return convertedAccess
}
