package accessAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/model"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) CreateAccess(ctx context.Context, req *desc.CreateRequest) (
	*desc.CreateResponse,
	error,
) {
	log.Printf("create access")

	if req.GetRoleId().Enum() == nil || req.GetPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Role and Path must be provided")
	}

	var createdAccess = &model.Access{
		RoleId: int(req.GetRoleId()) + 1,
		Path:   req.GetPath(),
	}

	id, err := i.accessService.CreateAccess(ctx, createdAccess)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
