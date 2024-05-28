package accessAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/logger"
	"github.com/semho/chat-microservices/auth/internal/model"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateAccess(ctx context.Context, req *desc.CreateRequest) (
	*desc.CreateResponse,
	error,
) {
	if req.GetRoleId() == 0 || req.GetPath() == "" {
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

	logger.Info("Created access with id:", zap.Int64("id", id))

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
