package main

import (
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"github.com/dupmanio/dupman/packages/user-api/app"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fxApp := fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return logWrapper.NewFxWrapper(logger)
		}),
		app.Provide(),
		fx.Invoke(app.Run),
	)

	fxApp.Run()
}
