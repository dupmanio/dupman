package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	Server    config.ServerConfig
	Database  config.DatabaseConfig
	Redis     config.RedisConfig
	Telemetry config.TelemetryConfig
	Keto      config.KetoConfig
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("notify", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
