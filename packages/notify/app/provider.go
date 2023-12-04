package app

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/database"
	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServer "github.com/dupmanio/dupman/packages/common/server"
	"github.com/dupmanio/dupman/packages/notify/config"
	"github.com/dupmanio/dupman/packages/notify/controller"
	"github.com/dupmanio/dupman/packages/notify/middleware"
	"github.com/dupmanio/dupman/packages/notify/repository"
	"github.com/dupmanio/dupman/packages/notify/route"
	"github.com/dupmanio/dupman/packages/notify/service"
	"github.com/dupmanio/dupman/packages/notify/version"
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
	ot, err := otel.NewOTel(
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
			fxHelper.AsRouteReceiver(serverProvider),
		),
		controller.Provide(),
		middleware.Provide(),
		route.Provide(),
		repository.Provide(),
		service.Provide(),
	)
}
