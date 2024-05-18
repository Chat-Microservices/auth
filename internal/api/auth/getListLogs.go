package authAPI

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/converter"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
)

func (i *Implementation) GetListLogs(ctx context.Context, req *desc.GetListLogsRequest) (*desc.LogsResponse, error) {
	pageNumber := req.GetPageNumber()
	pageSize := req.GetPageSize()

	if pageNumber == 0 {
		pageNumber = 1
	}

	if pageSize == 0 {
		pageSize = 5
	}

	listLogs, err := i.authService.GetListLogs(ctx, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	return converter.ToLogFromService(listLogs), nil
}
