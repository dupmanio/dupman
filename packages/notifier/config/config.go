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

type Mailer struct {
	Host     string `mapstructure:"MAILER_HOST"`
	Port     int    `mapstructure:"MAILER_PORT"`
	Username string `mapstructure:"MAILER_USERNAME"`
	Password string `mapstructure:"MAILER_PASSWORD"`
}

type Email struct {
	From string `mapstructure:"EMAIL_FROM"`
}

type Deliverer struct {
	Retries int `mapstructure:"DELIVERER_RETRIES" default:"3"`
}

type DupmanAPIService struct {
	ClientID     string   `mapstructure:"DUPMAN_API_SERVICE_CLIENT_ID"`
	ClientSecret string   `mapstructure:"DUPMAN_API_SERVICE_CLIENT_SECRET"`
	Scopes       []string `mapstructure:"DUPMAN_API_SERVICE_SCOPES" default:"[user:get_contact_info]"`
}

type Config struct {
	Env              string           `mapstructure:"ENV" default:"prod"`
	RabbitMQ         RabbitMQ         `mapstructure:",squash"`
	Worker           Worker           `mapstructure:",squash"`
	Mailer           Mailer           `mapstructure:",squash"`
	Deliverer        Deliverer        `mapstructure:",squash"`
	Email            Email            `mapstructure:",squash"`
	DupmanAPIService DupmanAPIService `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("packages/notifier", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
