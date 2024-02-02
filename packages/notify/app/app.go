package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/notify/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Server    *http.Server
	Migrators []fxHelper.IMigrator `group:"migrators"`
	Logger    *zap.Logger
	Config    *config.Config
	OT        *otel.OTel
}

func Run(params Params, lc fx.Lifecycle) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fxHelper.Migrate(params.Logger, params.Migrators...)

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

			// Shutdown does not close the open SSE connections,
			// so we have to use Close here.
			if err := params.Server.Close(); err != nil {
				return fmt.Errorf("failed to shutdown server: %w", err)
			}

			if err := params.OT.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown telemetry service: %w", err)
			}

			return nil
		},
	})

	return nil
}
