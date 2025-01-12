package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/service"
	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/vault"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Server       *http.Server
	Migrators    []fxHelper.IMigrator `group:"migrators"`
	Logger       *zap.Logger
	Config       *config.Config
	OT           *otel.OTel
	MessengerSvc *service.MessengerService
	Vault        *vault.Vault
}

func Run(params Params, lc fx.Lifecycle) error {
	vaultRenewerCtx, vaultRenewerCtxCancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			fxHelper.Migrate(params.Logger, params.Migrators...)

			params.Logger.Info("Starting Vault Renewer")

			go func() {
				if err := params.Vault.PeriodicallyRenewLeases(vaultRenewerCtx); err != nil {
					params.Logger.Fatal("Unable to start Vault Renewer", zap.Error(err))
				}
			}()

			params.Logger.Info(
				"Starting HTTP Server",
				zap.String(string(semconv.ServerAddressKey), params.Config.Server.ListenAddr),
				zap.String(string(semconv.ServerPortKey), params.Config.Server.Port),
			)

			go func() {
				if err := params.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					params.Logger.Fatal("Unable to start HTTP Server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Shutting down HTTP Server")

			if err := params.Server.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown server: %w", err)
			}

			if err := params.MessengerSvc.Close(); err != nil {
				return fmt.Errorf("failed to close messenger: %w", err)
			}

			vaultRenewerCtxCancel()

			if err := params.OT.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown telemetry service: %w", err)
			}

			return nil
		},
	})

	return nil
}
