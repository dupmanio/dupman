package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/api/route"
	"github.com/dupmanio/dupman/packages/api/server"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/common/logger"
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return logWrapper.NewFxWrapper(logger)
		}),
		fx.Provide(
			config.New,
			server.New,
			database.New,
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
		repository.Create(),
		service.Create(),
		fx.Invoke(server.Run),
	)

	app.Run()
}
