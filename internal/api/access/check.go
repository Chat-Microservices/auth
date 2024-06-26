package accessAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	logger.Info("endpoint:", zap.String("endpoint", req.GetEndpointAddress()))
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
