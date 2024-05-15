package converter

import (
	"github.com/semho/chat-microservices/auth/internal/model"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
)

func ToAccessFromService(list []*model.Access) *desc.ListAccessResponse {
	accessResponses := make([]*desc.Access, len(list))
	for i, log := range list {
		accessResponses[i] = &desc.Access{
			Id:     log.ID,
			RoleId: desc.Role(log.RoleId - 1),
			Path:   log.Path,
		}

	}
	return &desc.ListAccessResponse{
		List: accessResponses,
	}
}

func ToAuthUpdateAccessFromDesc(detail *desc.UpdateRequest) *model.Access {
	return &model.Access{
		ID:     detail.GetId(),
		RoleId: int(detail.GetRoleId().GetValue()) + 1,
		Path:   detail.GetPath().GetValue(),
	}
}
