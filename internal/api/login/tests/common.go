package tests

import "github.com/semho/chat-microservices/auth/internal/logger"

func initLogger() {
	if logger.Logger() == nil {
		logger.InitDefault("info")
	}
}
