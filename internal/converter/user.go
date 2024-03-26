package converter

import (
	"fmt"
	"github.com/semho/chat-microservices/auth/internal/model"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// из сервиса в протобаф
func ToAuthFromService(user *model.User) *desc.UserResponse {
	return &desc.UserResponse{
		Id:        user.ID,
		Detail:    ToAuthDetailFromService(user.Detail),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToAuthDetailFromService(detail model.Detail) *desc.UserDetail {
	return &desc.UserDetail{
		Name:  detail.Name,
		Email: detail.Email,
		Role:  desc.Role(detail.Role - 1).Enum(),
	}
}

// из протобаф в модель сервиса
func ToAuthDetailFromDesc(detail *desc.UserDetail) *model.Detail {
	return &model.Detail{
		Name:  detail.Name,
		Email: detail.Email,
		Role:  int(roleToInt(detail.GetRole())),
	}
}

func ToAuthUpdateUserFromDesc(detail *desc.UpdateRequest) *model.UpdateUserRequest {
	return &model.UpdateUserRequest{
		ID:    detail.GetId(),
		Name:  detail.GetInfo().GetName().GetValue(),
		Email: detail.GetInfo().GetEmail().GetValue(),
	}
}

// вспомогательные
func roleToInt(roleEnum desc.Role) int32 {
	role := fmt.Sprint(roleEnum)
	roleValue, ok := desc.Role_value[role]
	//костыль для записи в БД, т.к. enum c 0, а в БД с 1
	if !ok {
		roleValue = int32(desc.Role_user) + 1
	} else {
		roleValue++
	}

	return roleValue
}
