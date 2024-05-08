package config

import (
	"github.com/joho/godotenv"
	"time"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	Address() string
	IpAddress() string
}

type SwaggerConfig interface {
	Address() string
}

type TokenConfig interface {
	Prefix() string
	RefreshData() (string, time.Duration)
	AccessData() (string, time.Duration)
}
