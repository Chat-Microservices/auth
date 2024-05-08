package converter

import (
	modelRepo "github.com/semho/chat-microservices/auth/internal/repository/access/model"
)

func ToMapAccessFromRepo(access []modelRepo.Access) map[string]int {
	accessibleRoles := make(map[string]int)

	for _, onRow := range access {
		accessibleRoles[onRow.Path] = int(onRow.RoleId)
	}

	return accessibleRoles
}
