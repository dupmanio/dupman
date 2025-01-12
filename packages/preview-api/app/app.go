package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(server *http.Server, lc fx.Lifecycle, logger *zap.Logger, config *config.Config, ot *otel.OTel) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info(
				"Starting HTTP Server",
				zap.String(string(semconv.ServerAddressKey), config.Server.ListenAddr),
				zap.String(string(semconv.ServerPortKey), config.Server.Port),
			)

			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP Server")

			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown server: %w", err)
			}

			if err := ot.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown telemetry service: %w", err)
			}

			return nil
		},
	})

	return nil
}
