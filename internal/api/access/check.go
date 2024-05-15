package accessAPI

import (
	"context"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	log.Printf("endpoint: %s", req.GetEndpointAddress())
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
