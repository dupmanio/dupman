package app

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/vault"
	"github.com/dupmanio/dupman/packages/scanner/processor"
	"github.com/dupmanio/dupman/packages/scanner/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(
	lc fx.Lifecycle,
	logger *zap.Logger,
	processor *processor.Processor,
	messengerSvc *service.MessengerService,
	vault *vault.Vault,
) error {
	vaultRenewerCtx, vaultRenewerCtxCancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("Starting Vault Renewer")

			go func() {
				if err := vault.PeriodicallyRenewLeases(vaultRenewerCtx); err != nil {
					logger.Fatal("Unable to start Vault Renewer", zap.Error(err))
				}
			}()

			logger.Info("Starting Worker")

			go func() {
				if err := processor.Process(); err != nil {
					logger.Fatal("Unable to process", zap.Error(err))
				}
			}()

			logger.Info("Worker has been started. Waiting for messages")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down Worker")

			if err := messengerSvc.Close(); err != nil {
				return fmt.Errorf("failed to close messenger: %w", err)
			}

			vaultRenewerCtxCancel()

			return nil
		},
	})

	return nil
}
