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

type ChromeConfig struct {
	ResolutionX int `mapstructure:"CHROME_RESOLUTION_X" default:"1920"`
	ResolutionY int `mapstructure:"CHROME_RESOLUTION_Y" default:"1080"`
	Timeout     int `mapstructure:"CHROME_TIMEOUT" default:"10"`
}

type TelemetryConfig struct {
	Enabled      bool   `mapstructure:"TELEMETRY_ENABLED" default:"false"`
	CollectorURL string `mapstructure:"TELEMETRY_COLLECTOR_URL"`
}

type Config struct {
	Env       string          `mapstructure:"ENV" default:"prod"`
	AppName   string          `mapstructure:"APP_NAME" default:"preview-api"`
	LogPath   string          `mapstructure:"LOG_PATH" default:"/var/log/app.log"`
	Server    ServerConfig    `mapstructure:",squash"`
	Chrome    ChromeConfig    `mapstructure:",squash"`
	Telemetry TelemetryConfig `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("packages/preview-api", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
