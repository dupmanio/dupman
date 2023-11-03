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

type OAuthConfig struct {
	Issuer string `mapstructure:"OAUTH_ISSUER"`
}

type ChromeConfig struct {
	ResolutionX int `mapstructure:"CHROME_RESOLUTION_X" default:"1920"`
	ResolutionY int `mapstructure:"CHROME_RESOLUTION_Y" default:"1080"`
	Timeout     int `mapstructure:"CHROME_TIMEOUT" default:"10"`
}

type Config struct {
	Env    string       `mapstructure:"ENV" default:"prod"`
	Server ServerConfig `mapstructure:",squash"`
	OAuth  OAuthConfig  `mapstructure:",squash"`
	Chrome ChromeConfig `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("packages/preview-api", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
