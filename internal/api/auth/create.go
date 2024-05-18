package authAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetDetail() == nil || req.GetPassword() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Detail and Password must be provided")
	}

	if req.GetPassword().Password != req.GetPassword().PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "Password and Password Confirm do not match")
	}

	logger.Info("Create with username:", zap.String("username", req.GetDetail().GetName()))
	id, err := i.authService.Create(
		ctx,
		converter.ToAuthDetailFromDesc(req.GetDetail()),
		req.GetPassword().GetPassword(),
	)
	if err != nil {
		return nil, err
	}

	logger.Info("Created user with id:", zap.Int64("id", id))

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
