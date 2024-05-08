package accessService

import (
	"context"
	"errors"
	"fmt"
	"github.com/semho/chat-microservices/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type AccessCache struct {
	AccessibleRoles map[string]int
	ExpirationTime  time.Time
}

var cache AccessCache

func (s serv) Check(ctx context.Context, endpoint string) error {
	if endpoint == "" {
		return status.Error(codes.InvalidArgument, "Invalid request: username and password must be provided")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no metadata in context")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("no authorization header in context")
	}

	if !strings.HasPrefix(authHeader[0], s.tokenConfig.Prefix()) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], s.tokenConfig.Prefix())

	accessTokenSecretKey, _ := s.tokenConfig.AccessData()
	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return err
	}

	accessibleMap, err := s.AccessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[endpoint]
	if !ok {
		return nil
	}

	if role == claims.Role {
		return nil
	}

	return errors.New("access denied")
}

// TODO: Возможно стоит сделать кэширование через Redis
// TODO: подумать об инвалидации кэша при событиях
func (s serv) AccessibleRoles(ctx context.Context) (map[string]int, error) {
	if time.Now().After(cache.ExpirationTime) {
		fmt.Println("Обновляем кэш")
		accessibleMap, err := s.accessRepository.AccessibleRoles(ctx)
		if err != nil {
			return nil, err
		}

		// Обновляем кэш
		cache.AccessibleRoles = accessibleMap
		cache.ExpirationTime = time.Now().Add(5 * time.Minute)
	}

	return cache.AccessibleRoles, nil
}
