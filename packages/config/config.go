package config

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type TelemetryConfig struct {
	Enabled      bool   `mapstructure:"TELEMETRY_ENABLED" default:"false"`
	CollectorURL string `mapstructure:"TELEMETRY_COLLECTOR_URL"`
}

type BaseConfig struct {
	Env       string          `mapstructure:"ENV" default:"prod"`
	AppName   string          `mapstructure:"APP_NAME"`
	LogPath   string          `mapstructure:"LOG_PATH" default:"/var/log/app.log"`
	Telemetry TelemetryConfig `mapstructure:",squash"`
}

type Config interface {
	SetAppName(string)
}

func (conf *BaseConfig) SetAppName(appName string) {
	conf.AppName = appName
}

func Load(appName string, conf Config) error {
	conf.SetAppName(appName)

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
