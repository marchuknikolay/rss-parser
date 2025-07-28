package config

import (
	"fmt"
	"time"

	"github.com/joeshaw/envdecode"
)

type DBConfig struct {
	Host          string `env:"DB_HOST, required"`
	User          string `env:"DB_USER, required"`
	Password      string `env:"DB_PASSWORD, required"`
	Name          string `env:"DB_NAME, required"`
	HostPort      int    `env:"DB_HOST_PORT, required"`
	ContainerPort int    `env:"DB_CONTAINER_PORT, required"`
}

type ServerConfig struct {
	Port            int           `env:"SERVER_PORT, required"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT, required"`
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

func New() (*Config, error) {
	var config Config

	if err := envdecode.StrictDecode(&config); err != nil {
		return nil, fmt.Errorf("failed decoding .env file: %w", err)
	}

	return &config, nil
}
