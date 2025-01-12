package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Email struct {
	From string
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ   config.RabbitMQConfig
	Worker     config.WorkerConfig
	Dupman     config.DupmanConfig
	Telemetry  config.TelemetryConfig
	ServiceURL config.ServiceURLConfig `mapstructure:"service_url"`

	Mailer Mailer
	Email  Email
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("notifier", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
