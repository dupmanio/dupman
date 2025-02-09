package app

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/notifier/processor"
	"github.com/dupmanio/dupman/packages/notifier/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(
	lc fx.Lifecycle,
	logger *zap.Logger,
	processor *processor.Processor,
	messengerSvc *service.MessengerService,
) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("Starting Worker")

			go func() {
				if err := processor.Process(); err != nil {
					logger.Fatal("Unable to process", zap.Error(err))
				}
			}()

			logger.Info("Worker has been started. Waiting for messages")

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("Shutting down Worker")

			if err := messengerSvc.Close(); err != nil {
				return fmt.Errorf("failed to close messenger: %w", err)
			}

			return nil
		},
	})

	return nil
}
