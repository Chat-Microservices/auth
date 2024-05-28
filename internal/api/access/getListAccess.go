package accessAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	desc "github.com/semho/chat-microservices/auth/pkg/access_v1"
)

func (i *Implementation) GetListAccess(ctx context.Context, req *desc.GetListAccessRequest) (
	*desc.ListAccessResponse,
	error,
) {
	pageNumber := req.GetPageNumber()
	pageSize := req.GetPageSize()

	if pageNumber == 0 {
		pageNumber = 1
	}

	if pageSize == 0 {
		pageSize = 5
	}

	listAccess, err := i.accessService.GetListAccess(ctx, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}
	return converter.ToAccessFromService(listAccess), nil
}
