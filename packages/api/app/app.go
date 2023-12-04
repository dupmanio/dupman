package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/common/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(
	server *http.Server,
	lc fx.Lifecycle,
	logger *zap.Logger,
	config *config.Config,
	messengerSvc *service.MessengerService,
	ot *otel.OTel,
) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(
				"Starting HTTP Server",
				zap.String(string(semconv.ServerAddressKey), config.Server.ListenAddr),
				zap.String(string(semconv.ServerPortKey), config.Server.Port),
			)

			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("Unable to start HTTP Server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP Server")

			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown server: %w", err)
			}

			if err := messengerSvc.Close(); err != nil {
				return fmt.Errorf("failed to close messenger: %w", err)
			}

			if err := ot.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown telemetry service: %w", err)
			}

			return nil
		},
	})

	return nil
}
