package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	Server    config.ServerConfig    `mapstructure:",squash"`
	Database  config.DatabaseConfig  `mapstructure:",squash"`
	Redis     config.RedisConfig     `mapstructure:",squash"`
	Telemetry config.TelemetryConfig `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("notify", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
