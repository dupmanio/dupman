package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
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

func main() {
	app := fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(
			config.New,
			server.New,
			// Crete logger.
			func(conf *config.Config) (*zap.Logger, error) {
				loggerInst, err := logger.New(conf.Env, conf.AppName, "1.0.0", conf.LogPath)
				if err != nil {
					return nil, fmt.Errorf("unable to create logger: %w", err)
				}

				return loggerInst, nil
			},
		),
		controller.Create(),
		middleware.Create(),
		route.Create(),
		service.Create(),
		fx.Invoke(server.Run),
	)

	app.Run()
}
