package config

import (
	"fmt"

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

func NewDBConfig() (*DBConfig, error) {
	var config DBConfig

	if err := envdecode.StrictDecode(&config); err != nil {
		return nil, fmt.Errorf("failed decoding .env file: %w", err)
	}

	return &config, nil
}
