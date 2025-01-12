package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/config"
)

type ChromeConfig struct {
	RemoteURL   string `mapstructure:"remote_url" default:"ws://127.0.0.1:3000"`
	ResolutionX int    `mapstructure:"resolution_x" default:"1920"`
	ResolutionY int    `mapstructure:"resolution_y" default:"1080"`
	Timeout     int    `mapstructure:"timeout" default:"10"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	Server     config.ServerConfig
	Telemetry  config.TelemetryConfig
	ServiceURL config.ServiceURLConfig

	Chrome ChromeConfig
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("preview-api", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
