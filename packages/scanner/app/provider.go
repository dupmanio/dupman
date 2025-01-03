package app

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/vault"
	"github.com/dupmanio/dupman/packages/scanner/config"
	"github.com/dupmanio/dupman/packages/scanner/fetcher"
	"github.com/dupmanio/dupman/packages/scanner/processor"
	"github.com/dupmanio/dupman/packages/scanner/service"
	"github.com/dupmanio/dupman/packages/scanner/version"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func loggerProvider(conf *config.Config) (*zap.Logger, error) {
	loggerInst, err := logger.New(conf.Env, conf.AppName, version.Version)
	if err != nil {
		return nil, fmt.Errorf("unable to create logger: %w", err)
	}

	return loggerInst, nil
}

func oTelProvider(conf *config.Config, logger *zap.Logger) (*otel.OTel, error) {
	ctx := context.Background()

	ot, err := otel.NewOTel(
		ctx,
		conf.Env,
		conf.AppName,
		version.Version,
		conf.Telemetry.CollectorURL,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Telemetry service: %w", err)
	}

	return ot, nil
}

func vaultProvider(conf *config.Config, logger *zap.Logger, ot *otel.OTel) (*vault.Vault, error) {
	vaultConf := vault.Config{
		Address:  conf.Vault.ServerAddress,
		RoleID:   conf.Vault.RoleID,
		SecretID: conf.Vault.SecretID,
	}

	vaultInst, err := vault.New(vaultConf, logger, ot)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Vault: %w", err)
	}

	return vaultInst, nil
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			config.New,
			fetcher.New,
			processor.NewProcessor,
			loggerProvider,
			oTelProvider,
			vaultProvider,
		),
		service.Provide(),
	)
}
