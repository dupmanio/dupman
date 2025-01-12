package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	Server    config.ServerConfig
	Database  config.DatabaseConfig
	RabbitMQ  config.RabbitMQConfig
	Telemetry config.TelemetryConfig
	Vault     config.VaultConfig
	Keto      config.KetoConfig
	Notifier  config.Exchange
	Scanner   config.Exchange
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("api", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
