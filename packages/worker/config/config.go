package config

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/config"
)

type DupmanConfig struct {
	ClientID     string   `mapstructure:"DUPMAN_CLIENT_ID"`
	ClientSecret string   `mapstructure:"DUPMAN_CLIENT_SECRET"`
	Scopes       []string `mapstructure:"DUPMAN_SCOPES"`
}

type Config struct {
	Env    string       `mapstructure:"ENV" default:"prod"`
	Dupman DupmanConfig `mapstructure:",squash"`
}

func New() (*Config, error) {
	conf := new(Config)
	if err := config.Load("packages/worker", conf); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return conf, nil
}
