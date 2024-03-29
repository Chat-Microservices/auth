package authRepository

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

	tableName2 = "roles"

	idColumn2 = "id"
	idName2   = "name"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r repo) Get(ctx context.Context, id int64) (*model.User, error) {
	exists, err := r.userExists(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	query, args, err := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "user_repository.Get",
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

func (r repo) Create(ctx context.Context, detail *model.Detail, pass string) (int64, error) {

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(detail.Name, detail.Email, pass, detail.Role).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRow: query,
	}

	var userID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID); err != nil {
		log.Printf("failed to insert user into the database: %v", err)
		return 0, status.Error(codes.Internal, "Internal server error")
	}

	return userID, nil
}

func (r repo) Update(ctx context.Context, updateUser *model.UpdateUserRequest) error {
	exists, err := r.userExists(ctx, updateUser.ID)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return status.Error(codes.NotFound, "User not found")
	}

	query, values, err := buildUpdateQuery(*updateUser)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRow: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, values...)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("Обновлено строк: %d", rowCount)

	return nil
}

func (r repo) Delete(ctx context.Context, id int64) error {
	exists, err := r.userExists(ctx, id)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return status.Error(codes.NotFound, "User not found")
	}

	query, args, err := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRow: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("удалено строк: %d", rowCount)

	return nil
}

func (r *repo) userExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	q := db.Query{
		Name:     "user_repository.Exist",
		QueryRow: "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
	}
	err := r.db.DB().QueryRowContext(ctx, q, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func buildUpdateQuery(user model.UpdateUserRequest) (string, []any, error) {
	columns := make(map[string]interface{})

	if user.Name != "" {
		columns[nameColumn] = user.Name
	}

	if user.Email != "" {
		columns[emailColumn] = user.Email
	}

	query, args, err := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		SetMap(columns).
		Where(sq.Eq{idColumn: user.ID}).
		ToSql()

	return query, args, err
}
