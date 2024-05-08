package accessRepository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/repository"
	"github.com/semho/chat-microservices/auth/internal/repository/access/converter"
	modelRepo "github.com/semho/chat-microservices/auth/internal/repository/access/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

const (
	tableName = "access"

	idColumn     = "id"
	roleIdColumn = "role_id"
	pathColumn   = "path"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

func (r repo) AccessibleRoles(ctx context.Context) (map[string]int, error) {
	query, args, err := sq.Select(idColumn, roleIdColumn, pathColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "access_repository.Get",
		QueryRow: query,
	}

	var access []modelRepo.Access
	err = r.db.DB().ScanAllContext(ctx, &access, q, args...)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return converter.ToMapAccessFromRepo(access), nil
}
