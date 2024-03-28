package authAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetDetail() == nil || req.GetPassword() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Detail and Password must be provided")
	}

	if req.GetPassword().Password != req.GetPassword().PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "Password and Password Confirm do not match")
	}

	log.Printf("User name: %v", req.GetDetail().GetName())
	id, err := i.authService.Create(ctx, converter.ToAuthDetailFromDesc(req.GetDetail()), req.GetPassword().GetPassword())
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
