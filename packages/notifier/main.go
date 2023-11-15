package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"github.com/dupmanio/dupman/packages/notifier/deliverer"
	"github.com/dupmanio/dupman/packages/notifier/processor"
	"github.com/dupmanio/dupman/packages/notifier/service"
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
			processor.NewProcessor,
			// Crete logger.
			func(conf *config.Config) (*zap.Logger, error) {
				loggerInst, err := logger.New(conf.Env, conf.AppName, "1.0.0", conf.LogPath)
				if err != nil {
					return nil, fmt.Errorf("unable to create logger: %w", err)
				}

				return loggerInst, nil
			},
		),
		deliverer.Create(),
		service.Create(),
		fx.Invoke(worker.Run),
	)

	app.Run()
}
