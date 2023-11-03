package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/config"
)

type ServerConfig struct {
	ListenAddr     string   `mapstructure:"SERVER_ADDR" default:"0.0.0.0"`
	Port           string   `mapstructure:"SERVER_PORT" default:"8080"`
	TrustedProxies []string `mapstructure:"TRUSTED_PROXIES"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Database string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
}

type OAuthConfig struct {
	Issuer string `mapstructure:"OAUTH_ISSUER"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"RMQ_HOST" default:"127.0.0.1"`
	Port     string `mapstructure:"RMQ_PORT" default:"5672"`
	User     string `mapstructure:"RMQ_USER"`
	Password string `mapstructure:"RMQ_PASSWORD"`
}

type Notify struct {
	ExchangeName string `mapstructure:"NOTIFY_EXCHANGE_NAME"`
	RoutingKey   string `mapstructure:"NOTIFY_ROUTING_KEY"`
}

type Config struct {
	Env      string         `mapstructure:"ENV" default:"prod"`
	Server   ServerConfig   `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	OAuth    OAuthConfig    `mapstructure:",squash"`
	RabbitMQ RabbitMQ       `mapstructure:",squash"`
	Notify   Notify         `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("packages/api", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
