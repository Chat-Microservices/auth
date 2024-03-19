package main

import (
	"context"
	"flag"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/semho/chat-microservices/auth/internal/config"
	"github.com/semho/chat-microservices/auth/internal/config/env"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

type UserDB struct {
	ID        int64
	Name      string
	Email     string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	ID    int64
	Name  string
	Email string
}

func (s *server) userExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func buildUpdateQuery(req UpdateUserRequest) (string, []any, error) {
	columns := make(map[string]interface{})

	if req.Name != "" {
		columns["name"] = req.Name
	}

	if req.Email != "" {
		columns["email"] = req.Email
	}

	query, args, err := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		SetMap(columns).
		Where(sq.Eq{"id": req.ID}).
		ToSql()

	return query, args, err
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	exists, err := s.userExists(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	query, args, err := sq.Select("id, name, email, role, created_at, updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	var user UserDB
	err = s.pool.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &desc.GetResponse{
		User: &desc.UserResponse{
			Id: req.GetId(),
			Detail: &desc.UserDetail{
				Name:  user.Name,
				Email: user.Email,
				Role:  desc.Role(user.Role - 1).Enum(),
			},
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetDetail() == nil || req.GetPassword() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Detail and Password must be provided")
	}

	if req.GetPassword().Password != req.GetPassword().PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "Password and Password Confirm do not match")
	}

	log.Printf("User name: %v", req.GetDetail().Name)
	//приводим к строке, иначе не верно будет браться id в enum
	role := fmt.Sprint(req.GetDetail().GetRole())
	roleValue, ok := desc.Role_value[role]
	//костыль для записи в БД, т.к. enum c 0, а в БД с 1
	if !ok {
		roleValue = int32(desc.Role_user) + 1
	} else {
		roleValue++
	}

	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name, email, password, role").
		Values(req.GetDetail().Name, req.GetDetail().Email, req.GetPassword().Password, roleValue).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	var userID int64
	if err = s.pool.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		log.Printf("failed to insert user into the database: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	newUser := &desc.User{
		Id:       userID,
		Detail:   req.GetDetail(),
		Password: req.GetPassword(),
	}
	log.Printf("user %v", newUser)

	response := &desc.CreateResponse{
		Id: newUser.Id,
	}

	return response, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	exists, err := s.userExists(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	updateRequest := UpdateUserRequest{
		ID:    req.GetId(),
		Name:  req.GetInfo().GetName().GetValue(),
		Email: req.GetInfo().GetEmail().GetValue(),
	}
	query, values, err := buildUpdateQuery(updateRequest)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	res, err := s.pool.Exec(ctx, query, values...)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("Обновлено строк: %d", rowCount)

	updateUser := &desc.UpdateUserInfo{
		Name:  req.GetInfo().GetName(),
		Email: req.GetInfo().GetEmail(),
	}

	log.Printf("update: %v", updateUser)

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	exists, err := s.userExists(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	query, args, err := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("удалено строк: %d", rowCount)

	log.Printf("delete user by id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
