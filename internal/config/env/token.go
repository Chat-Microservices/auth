package env

import (
	"errors"
	"github.com/semho/chat-microservices/auth/internal/config"
	"os"
	"time"
)

const (
	AuthPrefix            = "AUTH_PREFIX"
	RefreshTokenSecretKey = "REFRESH_TOKEN_KEY"
	AccessTokenSecretKey  = "ACCESS_TOKEN_KEY"

	RefreshTokenExpiration = "REFRESH_TOKEN_EXPIRATION"
	AccessTokenExpiration  = "ACCESS_TOKEN_EXPIRATION"
)

type TokenConfig struct {
	prefix            string
	refreshSecret     string
	accessSecret      string
	refreshExpiration time.Duration
	accessExpiration  time.Duration
}

func NewTokenConfig() (config.TokenConfig, error) {
	prefix := os.Getenv(AuthPrefix)
	if len(prefix) == 0 {
		return nil, errors.New("AuthPrefix not found")
	}
	refreshSecret := os.Getenv(RefreshTokenSecretKey)
	if len(refreshSecret) == 0 {
		return nil, errors.New("RefreshTokenSecretKey not found")
	}
	accessSecret := os.Getenv(AccessTokenSecretKey)
	if len(accessSecret) == 0 {
		return nil, errors.New("AccessTokenSecretKey not found")
	}

	refreshExpiration := os.Getenv(RefreshTokenExpiration)
	if len(refreshExpiration) == 0 {
		return nil, errors.New("RefreshTokenExpiration not found")
	}
	refreshExpirationDuration, err := time.ParseDuration(refreshExpiration)
	if err != nil {
		return nil, errors.New("error parsing refresh token expiration")
	}

	accessExpiration := os.Getenv(AccessTokenExpiration)
	if len(accessExpiration) == 0 {
		return nil, errors.New("AccessTokenExpiration not found")
	}
	accessExpirationDuration, err := time.ParseDuration(accessExpiration)
	if err != nil {
		return nil, errors.New("error parsing access token expiration")
	}

	return &TokenConfig{
		prefix:            prefix,
		refreshSecret:     refreshSecret,
		accessSecret:      accessSecret,
		refreshExpiration: refreshExpirationDuration,
		accessExpiration:  accessExpirationDuration,
	}, nil
}

func (t *TokenConfig) Prefix() string {
	return t.prefix
}

func (t *TokenConfig) RefreshData() (string, time.Duration) {
	return t.refreshSecret, t.refreshExpiration
}

func (t *TokenConfig) AccessData() (string, time.Duration) {
	return t.accessSecret, t.accessExpiration
}
