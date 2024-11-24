package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Notifier struct {
	ExchangeName string `mapstructure:"NOTIFIER_EXCHANGE_NAME"`
	RoutingKey   string `mapstructure:"NOTIFIER_ROUTING_KEY"`
}

type Scanner struct {
	ExchangeName string `mapstructure:"SCANNER_EXCHANGE_NAME"`
	RoutingKey   string `mapstructure:"SCANNER_ROUTING_KEY"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	Server    config.ServerConfig    `mapstructure:",squash"`
	Database  config.DatabaseConfig  `mapstructure:",squash"`
	RabbitMQ  config.RabbitMQConfig  `mapstructure:",squash"`
	Telemetry config.TelemetryConfig `mapstructure:",squash"`
	Vault     config.VaultConfig     `mapstructure:",squash"`
	Keto      config.KetoConfig      `mapstructure:",squash"`

	Notifier Notifier `mapstructure:",squash"`
	Scanner  Scanner  `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("api", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
