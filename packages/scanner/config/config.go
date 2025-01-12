package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ   config.RabbitMQConfig
	Worker     config.WorkerConfig
	Dupman     config.DupmanConfig
	Telemetry  config.TelemetryConfig
	Vault      config.VaultConfig
	ServiceURL config.ServiceURLConfig
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("scanner", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
