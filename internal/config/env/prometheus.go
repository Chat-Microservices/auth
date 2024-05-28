package env

import (
	"errors"
	"github.com/semho/chat-microservices/auth/internal/config"
	"net"
	"os"
)

const (
	prometheusHTTPHostEnvName = "PROMETHEUS_HTTP_HOST"
	prometheusHTTPPortEnvName = "PROMETHEUS_HTTP_PORT"
)

type prometheusConfig struct {
	hostHTTP string
	portHTTP string
}

func NewPrometheusConfig() (config.PrometheusConfig, error) {
	hostHTTP := os.Getenv(prometheusHTTPHostEnvName)
	if len(hostHTTP) == 0 {
		return nil, errors.New("prometheus http host not found")
	}
	portHTTP := os.Getenv(prometheusHTTPPortEnvName)
	if len(portHTTP) == 0 {
		return nil, errors.New("prometheus http port not found")
	}

	return &prometheusConfig{
		hostHTTP: hostHTTP,
		portHTTP: portHTTP,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.hostHTTP, cfg.portHTTP)
}
