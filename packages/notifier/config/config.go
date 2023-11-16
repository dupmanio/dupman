package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type Mailer struct {
	Host     string `mapstructure:"MAILER_HOST"`
	Port     int    `mapstructure:"MAILER_PORT"`
	Username string `mapstructure:"MAILER_USERNAME"`
	Password string `mapstructure:"MAILER_PASSWORD"`
}

type Email struct {
	From string `mapstructure:"EMAIL_FROM"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ config.RabbitMQConfig `mapstructure:",squash"`
	Worker   config.WorkerConfig   `mapstructure:",squash"`
	Dupman   config.DupmanConfig   `mapstructure:",squash"`

	Mailer Mailer `mapstructure:",squash"`
	Email  Email  `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("notifier", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
