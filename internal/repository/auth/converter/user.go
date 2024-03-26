package converter

import (
	"github.com/semho/chat-microservices/auth/internal/model"
	modelRepo "github.com/semho/chat-microservices/auth/internal/repository/auth/model"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Detail:    ToUserDetailFromRepo(user.Detail),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserDetailFromRepo(detail modelRepo.Detail) model.Detail {
	return model.Detail{
		Name:  detail.Name,
		Email: detail.Email,
		Role:  detail.Role,
	}
}
