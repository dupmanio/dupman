package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/config"
)

type DupmanConfig struct {
	ClientID     string `mapstructure:"DUPMAN_CLIENT_ID"`
	ClientSecret string `mapstructure:"DUPMAN_CLIENT_SECRET"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"RMQ_HOST" default:"127.0.0.1"`
	Port     string `mapstructure:"RMQ_PORT" default:"5672"`
	User     string `mapstructure:"RMQ_USER"`
	Password string `mapstructure:"RMQ_PASSWORD"`
}

type Scanner struct {
	ExchangeName string `mapstructure:"SCANNER_EXCHANGE_NAME"`
	RoutingKey   string `mapstructure:"SCANNER_ROUTING_KEY"`
}

type Config struct {
	Env      string       `mapstructure:"ENV" default:"prod"`
	Dupman   DupmanConfig `mapstructure:",squash"`
	RabbitMQ RabbitMQ     `mapstructure:",squash"`
	Scanner  Scanner      `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("packages/scanner-scheduler", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
