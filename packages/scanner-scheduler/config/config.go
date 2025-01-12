package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ   config.RabbitMQConfig
	Dupman     config.DupmanConfig
	Telemetry  config.TelemetryConfig
	ServiceURL config.ServiceURLConfig
	Scanner    config.Exchange
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("scanner-scheduler", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
