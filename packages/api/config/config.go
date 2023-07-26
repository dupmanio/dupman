package config

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
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

type CORSConfig struct {
	AllowOrigins  []string `mapstructure:"CORS_ALLOW_ORIGINS" default:"[*]"`
	AllowMethods  []string `mapstructure:"CORS_ALLOW_METHODS" default:"[*]"`
	AllowHeaders  []string `mapstructure:"CORS_ALLOW_HEADERS" default:"[*]"`
	ExposeHeaders []string `mapstructure:"CORS_EXPOSE_HEADERS" default:"[*]"`
}

type OAuthConfig struct {
	Issuer string `mapstructure:"OAUTH_ISSUER"`
}

type Config struct {
	Env      string         `mapstructure:"ENV" default:"prod"`
	Server   ServerConfig   `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	CORS     CORSConfig     `mapstructure:",squash"`
	OAuth    OAuthConfig    `mapstructure:",squash"`
}

func New() (*Config, error) {
	// @todo: refactor.
	conf := new(Config)

	viper.AddConfigPath("packages/api")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	defaults.SetDefaults(conf)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w	", err)
	}

	return conf, nil
}
