package accessRepository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/model"
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

func (r repo) GetListAccess(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Access, error) {
	offset := (pageNumber - 1) * pageSize

	query, args, err := sq.Select(idColumn, roleIdColumn, pathColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Limit(pageSize).
		Offset(offset).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "access_repository.GetList",
		QueryRow: query,
	}

	var listAccess []modelRepo.Access
	err = r.db.DB().ScanAllContext(ctx, &listAccess, q, args...)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return converter.ToAccessListFromRepo(listAccess), nil
}

func (r repo) CreateAccess(ctx context.Context, roleId int, path string) (int64, error) {

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(roleIdColumn, pathColumn).
		Values(roleId, path).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "access_repository.Create",
		QueryRow: query,
	}

	var accessID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&accessID); err != nil {
		log.Printf("failed to insert access into the database: %v", err)
		return 0, status.Error(codes.Internal, "Internal server error")
	}

	return accessID, nil
}

func (r repo) DeleteAccess(ctx context.Context, id int64) error {
	exists, err := r.accessExists(ctx, id)
	if err != nil {
		return status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return status.Error(codes.NotFound, "access record not found")
	}

	query, args, err := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "access_repository.Delete",
		QueryRow: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("удалено строк: %d", rowCount)

	return nil
}

func (r *repo) accessExists(ctx context.Context, accessID int64) (bool, error) {
	var exists bool
	q := db.Query{
		Name:     "access_repository.Exist",
		QueryRow: "SELECT EXISTS(SELECT 1 FROM access WHERE id = $1)",
	}
	err := r.db.DB().QueryRowContext(ctx, q, accessID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r repo) UpdateAccess(ctx context.Context, access *model.Access) error {
	exists, err := r.accessExists(ctx, access.ID)
	if err != nil {
		return status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return status.Error(codes.NotFound, "Access record not found")
	}

	query, values, err := buildUpdateQuery(*access)
	if err != nil {
		return status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "access_repository.Update",
		QueryRow: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, values...)
	if err != nil {
		return status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("Обновлено строк: %d", rowCount)

	return nil
}

func buildUpdateQuery(access model.Access) (string, []any, error) {
	columns := make(map[string]interface{})

	if access.RoleId != 0 {
		columns[roleIdColumn] = access.RoleId
	}

	if access.Path != "" {
		columns[pathColumn] = access.Path
	}

	query, args, err := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		SetMap(columns).
		Where(sq.Eq{idColumn: access.ID}).
		ToSql()

	return query, args, err
}
