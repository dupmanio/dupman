package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ config.RabbitMQConfig `mapstructure:",squash"`
	Worker   config.WorkerConfig   `mapstructure:",squash"`
	Dupman   config.DupmanConfig   `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("scanner", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
