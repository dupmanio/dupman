package main

import (
	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/api/route"
	"github.com/dupmanio/dupman/packages/api/server"
	"github.com/dupmanio/dupman/packages/api/service"
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
			zap.NewProduction,
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
