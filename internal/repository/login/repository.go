package loginRepository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
	"github.com/semho/chat-microservices/auth/internal/repository"
	"github.com/semho/chat-microservices/auth/internal/repository/auth/converter"
	modelRepo "github.com/semho/chat-microservices/auth/internal/repository/auth/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.LoginRepository {
	return &repo{db: db}
}

func (r repo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query, args, err := sq.Select(
		idColumn,
		nameColumn,
		emailColumn,
		roleColumn,
		passwordColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: username}).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "login_repository.Get",
		QueryRow: query,
	}
	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)

	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return converter.ToUserFromRepo(&user), nil
}
