package accessAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) UpdateAccess(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}
	logger.Info("Updating access with id:", zap.Int64("id", req.GetId()))
	err := i.accessService.UpdateAccess(ctx, converter.ToAuthUpdateAccessFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
