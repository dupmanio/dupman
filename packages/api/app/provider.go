package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/migrator"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/api/route"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/api/version"
	"github.com/dupmanio/dupman/packages/common/database"
	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/common/ory/keto"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServer "github.com/dupmanio/dupman/packages/common/server"
	"github.com/dupmanio/dupman/packages/common/vault"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func loggerProvider(conf *config.Config) (*zap.Logger, error) {
	loggerInst, err := logger.New(conf.Env, conf.AppName, version.Version, conf.LogPath)
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

func databaseProvider(conf *config.Config, logger *zap.Logger, ot *otel.OTel) (*database.Database, error) {
	db, err := database.New(conf.Database, logger, ot)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return db, nil
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

func ketoProvider(conf *config.Config, logger *zap.Logger, ot *otel.OTel) (*keto.Keto, error) {
	ketoConf := keto.Config{
		WriteURL: conf.Keto.WriteURL,
	}

	ketoInst, err := keto.New(ketoConf, logger, ot)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Keto: %w", err)
	}

	return ketoInst, nil
}

func serverProvider(
	routes []fxHelper.IRoute,
	logger *zap.Logger,
	conf *config.Config,
	ot *otel.OTel,
) (*http.Server, error) {
	srv, err := commonServer.New(conf.Env, conf.Server, routes, logger, ot)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return srv, nil
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			config.New,
			loggerProvider,
			oTelProvider,
			databaseProvider,
			vaultProvider,
			ketoProvider,
			fxHelper.AsRouteReceiver(serverProvider),
		),
		controller.Provide(),
		migrator.Provide(),
		repository.Provide(),
		route.Provide(),
		service.Provide(),
	)
}
