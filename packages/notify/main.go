package main

import (
	"fmt"

	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"github.com/dupmanio/dupman/packages/common/logger"
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/notify/config"
	"github.com/dupmanio/dupman/packages/notify/controller"
	"github.com/dupmanio/dupman/packages/notify/database"
	"github.com/dupmanio/dupman/packages/notify/middleware"
	"github.com/dupmanio/dupman/packages/notify/repository"
	"github.com/dupmanio/dupman/packages/notify/route"
	"github.com/dupmanio/dupman/packages/notify/server"
	"github.com/dupmanio/dupman/packages/notify/service"
	"github.com/dupmanio/dupman/packages/notify/version"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// @todo: move provider functions from main package.

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

func main() {
	app := fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return logWrapper.NewFxWrapper(logger)
		}),
		fx.Provide(
			config.New,
			fxHelper.AsRouteReceiver(server.New),
			database.New,
			loggerProvider,
			oTelProvider,
		),
		controller.Provide(),
		middleware.Provide(),
		route.Provide(),
		repository.Provide(),
		service.Provide(),
		fx.Invoke(server.Run),
	)

	app.Run()
}
