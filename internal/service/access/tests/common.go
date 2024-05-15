package tests

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"time"
)

type mockTxManager struct{}

func (m *mockTxManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}

type MockTokenConfig struct{}

func (m *MockTokenConfig) Prefix() string {
	return "Bearer "
}

func (m *MockTokenConfig) RefreshData() (string, time.Duration) {
	return "refreshToken", time.Hour
}

func (m *MockTokenConfig) AccessData() (string, time.Duration) {
	return "accessToken", time.Minute * 15
}
