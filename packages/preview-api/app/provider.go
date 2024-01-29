package app

import (
	"context"
	"fmt"
	"net/http"

	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServer "github.com/dupmanio/dupman/packages/common/server"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	"github.com/dupmanio/dupman/packages/preview-api/controller"
	"github.com/dupmanio/dupman/packages/preview-api/route"
	"github.com/dupmanio/dupman/packages/preview-api/service"
	"github.com/dupmanio/dupman/packages/preview-api/version"
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
			fxHelper.AsRouteReceiver(serverProvider),
		),
		controller.Provide(),
		route.Provide(),
		service.Provide(),
	)
}
