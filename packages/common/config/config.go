package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	Env     string `validate:"oneof=dev test prod" default:"prod"`
	AppName string
}

type TelemetryConfig struct {
	CollectorURL string `mapstructure:"collector_url"`
}

type ServerConfig struct {
	ListenAddr     string   `mapstructure:"address" default:"0.0.0.0"`
	Port           string   `mapstructure:"port" default:"8080"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     string
}

type RedisConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RabbitMQConfig struct {
	Host     string `default:"127.0.0.1"`
	Port     string `default:"5672"`
	User     string
	Password string
}

type WorkerConfig struct {
	QueueName     string `mapstructure:"queue_name"`
	PrefetchCount int    `mapstructure:"prefetch_count" default:"1"`
	PrefetchSize  int    `mapstructure:"prefetch_size" default:"0"`
	RetryAttempts int    `mapstructure:"retry_attempts" default:"3"`
}

type DupmanConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Scopes       []string
	Audience     []string
}

type VaultConfig struct {
	ServerAddress string `mapstructure:"address"`
	RoleID        string `mapstructure:"role_id"`
	SecretID      string `mapstructure:"secret_id"`
}

type KetoConfig struct {
	WriteURL string `mapstructure:"write_url"`
}

type Exchange struct {
	ExchangeName string `mapstructure:"exchange_name"`
	RoutingKey   string `mapstructure:"routing_key"`
}

type ServiceURLConfig struct {
	// @todo: update URLs.
	API        string `mapstructure:"api" default:"http://gateway.dupman.localhost/api"`
	UserAPI    string `mapstructure:"user_api" default:"http://gateway.dupman.localhost/user-api"`
	PreviewAPI string `mapstructure:"preview_api" default:"http://gateway.dupman.localhost/preview-api"`
	Notify     string `mapstructure:"notify" default:"http://gateway.dupman.localhost/notify"`
}

type Config interface {
	SetAppName(appName string)
}

func (conf *BaseConfig) SetAppName(appName string) {
	conf.AppName = appName
}

func Load(appName string, conf Config) error {
	conf.SetAppName(appName)

	viperInstance := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
	)

	// @todo: add config path from CLI argument
	viperInstance.AddConfigPath(".")
	viperInstance.AddConfigPath(fmt.Sprintf("packages/%s", appName))
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yaml")

	viperInstance.AllowEmptyEnv(true)
	viperInstance.AutomaticEnv()

	defaults.SetDefaults(conf)

	if err := viperInstance.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viperInstance.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}
