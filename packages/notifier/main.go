package main

import (
	"github.com/dupmanio/dupman/packages/notifier/broker"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"github.com/dupmanio/dupman/packages/notifier/deliverer"
	"github.com/dupmanio/dupman/packages/notifier/processor"
	"github.com/dupmanio/dupman/packages/notifier/worker"
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
			zap.NewDevelopment,
			broker.NewRabbitMQ,
			processor.NewProcessor,
		),
		deliverer.Create(),
		fx.Invoke(worker.Run),
	)

	app.Run()
}
