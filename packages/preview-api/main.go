package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	"github.com/dupmanio/dupman/packages/preview-api/controller"
	"github.com/dupmanio/dupman/packages/preview-api/middleware"
	"github.com/dupmanio/dupman/packages/preview-api/route"
	"github.com/dupmanio/dupman/packages/preview-api/server"
	"github.com/dupmanio/dupman/packages/preview-api/service"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// @todo: move provider functions from main package.

func loggerProvider(conf *config.Config) (*zap.Logger, error) {
	// @todo: store app version separately.
	loggerInst, err := logger.New(conf.Env, conf.AppName, "1.0.0", conf.LogPath)
	if err != nil {
		return nil, fmt.Errorf("unable to create logger: %w", err)
	}

	return loggerInst, nil
}

func oTelProvider(conf *config.Config, logger *zap.Logger) (*otel.OTel, error) {
	ot, err := otel.NewOTel(
		conf.Env,
		conf.AppName,
		"1.0.0",
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
			server.New,
			loggerProvider,
			oTelProvider,
		),
		controller.Create(),
		middleware.Create(),
		route.Create(),
		service.Create(),
		fx.Invoke(server.Run),
	)

	app.Run()
}
