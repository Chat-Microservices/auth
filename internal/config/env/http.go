package env

import (
	"errors"
	"github.com/semho/chat-microservices/auth/internal/config"
	"net"
	"os"
)

var _ config.HTTPConfig = (*httpConfig)(nil)

const (
	httpHostEnvName   = "HTTP_HOST"
	httpPortEnvName   = "HTTP_PORT"
	httpIPHostEnvName = "HTTP_IP_HOST"
)

type httpConfig struct {
	host   string
	port   string
	ipHost string
}

func NewHTTPConfig() (*httpConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}
	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}
	ipHost := os.Getenv(httpIPHostEnvName)
	if len(ipHost) == 0 {
		ipHost = host
	}

	return &httpConfig{
		host:   host,
		port:   port,
		ipHost: ipHost,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *httpConfig) IpAddress() string {
	return net.JoinHostPort(cfg.ipHost, cfg.port)
}
