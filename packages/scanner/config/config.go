package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/config"
)

type RabbitMQ struct {
	Host      string `mapstructure:"RMQ_HOST" default:"127.0.0.1"`
	Port      string `mapstructure:"RMQ_PORT" default:"5672"`
	User      string `mapstructure:"RMQ_USER"`
	Password  string `mapstructure:"RMQ_PASSWORD"`
	QueueName string `mapstructure:"RMQ_QUEUE_NAME"`
}

type Worker struct {
	PrefetchCount int `mapstructure:"WORKER_PREFETCH_COUNT" default:"1"`
	PrefetchSize  int `mapstructure:"WORKER_PREFETCH_SIZE" default:"0"`
	RetryAttempts int `mapstructure:"WORKER_RETRY_ATTEMPTS" default:"3"`
}

type DupmanConfig struct {
	ClientID     string `mapstructure:"DUPMAN_CLIENT_ID"`
	ClientSecret string `mapstructure:"DUPMAN_CLIENT_SECRET"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	RabbitMQ     RabbitMQ     `mapstructure:",squash"`
	Worker       Worker       `mapstructure:",squash"`
	DupmanConfig DupmanConfig `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("scanner", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
