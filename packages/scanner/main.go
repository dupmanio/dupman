package main

import (
	"github.com/dupmanio/dupman/packages/scanner/broker"
	"github.com/dupmanio/dupman/packages/scanner/config"
	"github.com/dupmanio/dupman/packages/scanner/fetcher"
	"github.com/dupmanio/dupman/packages/scanner/processor"
	"github.com/dupmanio/dupman/packages/scanner/service"
	"github.com/dupmanio/dupman/packages/scanner/worker"
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
			fetcher.New,
			processor.NewProcessor,
		),
		service.Create(),
		fx.Invoke(worker.Run),
	)

	app.Run()
}
