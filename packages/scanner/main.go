package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
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
			fetcher.New,
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
		service.Create(),
		fx.Invoke(worker.Run),
	)

	app.Run()
}
