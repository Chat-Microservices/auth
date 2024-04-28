package env

import (
	"errors"
	"github.com/semho/chat-microservices/auth/internal/config"
	"net"
	"os"
)

var _ config.HTTPConfig = (*httpConfig)(nil)

const (
	httpHostEnvName = "HTTP_HOST"
	HttpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (*httpConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}
	port := os.Getenv(HttpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
