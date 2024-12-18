package config

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	Env     string `mapstructure:"ENV" default:"prod"`
	AppName string `mapstructure:"APP_NAME"`
	LogPath string `mapstructure:"LOG_PATH" default:"/var/log/app.log"`
}

type TelemetryConfig struct {
	CollectorURL string `mapstructure:"TELEMETRY_COLLECTOR_URL"`
}

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

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	User     string `mapstructure:"REDIS_USER"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	Database string `mapstructure:"REDIS_DB"`
	Port     string `mapstructure:"REDIS_PORT"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"RMQ_HOST" default:"127.0.0.1"`
	Port     string `mapstructure:"RMQ_PORT" default:"5672"`
	User     string `mapstructure:"RMQ_USER"`
	Password string `mapstructure:"RMQ_PASSWORD"`
}

type WorkerConfig struct {
	QueueName     string `mapstructure:"WORKER_QUEUE_NAME"`
	PrefetchCount int    `mapstructure:"WORKER_PREFETCH_COUNT" default:"1"`
	PrefetchSize  int    `mapstructure:"WORKER_PREFETCH_SIZE" default:"0"`
	RetryAttempts int    `mapstructure:"WORKER_RETRY_ATTEMPTS" default:"3"`
}

type DupmanConfig struct {
	ClientID     string   `mapstructure:"DUPMAN_CLIENT_ID"`
	ClientSecret string   `mapstructure:"DUPMAN_CLIENT_SECRET"`
	Scopes       []string `mapstructure:"DUPMAN_CLIENT_SCOPES"`
	Audience     []string `mapstructure:"DUPMAN_CLIENT_AUDIENCE"`
}

type VaultConfig struct {
	ServerAddress string `mapstructure:"VAULT_SERVER_ADDRESS"`
	RoleID        string `mapstructure:"VAULT_ROLE_ID"`
	SecretID      string `mapstructure:"VAULT_SECRET_ID"`
}

type KetoConfig struct {
	WriteURL string `mapstructure:"KETO_WRITE_URL"`
}

type ServiceURLConfig struct {
	// @todo: update URLs.
	API        string `mapstructure:"SERVICE_API_URL" default:"http://gateway.dupman.localhost/api"`
	UserAPI    string `mapstructure:"SERVICE_USER_API_URL" default:"http://gateway.dupman.localhost/user-api"`
	PreviewAPI string `mapstructure:"SERVICE_PREVIEW_API_URL" default:"http://gateway.dupman.localhost/preview-api"`
	Notify     string `mapstructure:"SERVICE_NOTIFY_URL" default:"http://gateway.dupman.localhost/notify"`
}

type Config interface {
	SetAppName(appName string)
}

func (conf *BaseConfig) SetAppName(appName string) {
	conf.AppName = appName
}

func Load(appName string, conf Config) error {
	// @todo: move to YAML format.
	conf.SetAppName(appName)

	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("packages/%s", appName))
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	defaults.SetDefaults(conf)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
