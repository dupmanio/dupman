package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/notify/config"
	"github.com/dupmanio/dupman/packages/notify/controller"
	"github.com/dupmanio/dupman/packages/notify/database"
	"github.com/dupmanio/dupman/packages/notify/middleware"
	"github.com/dupmanio/dupman/packages/notify/repository"
	"github.com/dupmanio/dupman/packages/notify/route"
	"github.com/dupmanio/dupman/packages/notify/server"
	"github.com/dupmanio/dupman/packages/notify/service"
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
