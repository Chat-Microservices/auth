package authAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	"github.com/semho/chat-microservices/auth/internal/logger"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"go.uber.org/zap"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.UserResponse, error) {
	logger.Info("Get user id", zap.Int64("id", req.GetId()))
	userObj, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToAuthFromService(userObj), nil
}
