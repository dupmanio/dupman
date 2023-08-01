package config

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
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
	// @todo: refactor.
	conf := new(Config)

	viper.AddConfigPath("packages/worker")
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
