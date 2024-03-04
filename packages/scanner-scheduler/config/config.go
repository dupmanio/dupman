package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Scanner struct {
	ExchangeName string `mapstructure:"SCANNER_EXCHANGE_NAME"`
	RoutingKey   string `mapstructure:"SCANNER_ROUTING_KEY"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ   config.RabbitMQConfig   `mapstructure:",squash"`
	Dupman     config.DupmanConfig     `mapstructure:",squash"`
	Telemetry  config.TelemetryConfig  `mapstructure:",squash"`
	ServiceURL config.ServiceURLConfig `mapstructure:",squash"`

	Scanner Scanner `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("scanner-scheduler", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
