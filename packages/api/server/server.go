package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/service"
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/gin-gonic/gin"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	Engine *gin.Engine
}

func New(logger *zap.Logger, config *config.Config, ot *otel.OTel) (*Server, error) {
	ginLogWrapper := logWrapper.NewGinWrapper(logger)

	gin.DefaultWriter = ginLogWrapper

	if config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.Env == "test" {
		gin.SetMode(gin.TestMode)
	}

	engine := gin.New()
	engine.ContextWithFallback = true

	if err := engine.SetTrustedProxies(config.Server.TrustedProxies); err != nil {
		return nil, fmt.Errorf("unable to set trusted proxies: %w", err)
	}

	engine.Use(ot.GetOTelGinMiddleware())
	engine.Use(ginLogWrapper.GetGinzapMiddleware())
	engine.Use(ginLogWrapper.GetGinzapRecoveryMiddleware())

	return &Server{
		Engine: engine,
	}, nil
}

func Run(
	server *Server,
	lc fx.Lifecycle,
	logger *zap.Logger,
	config *config.Config,
	messengerSvc *service.MessengerService,
	ot *otel.OTel,
) error {
	httpServer := http.Server{
		Addr:              net.JoinHostPort(config.Server.ListenAddr, config.Server.Port),
		Handler:           server.Engine,
		ReadHeaderTimeout: 0,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(
				"Starting HTTP Server",
				zap.String(string(semconv.ServerAddressKey), config.Server.ListenAddr),
				zap.String(string(semconv.ServerPortKey), config.Server.Port),
			)

			go func() {
				if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP Server")

			if err := httpServer.Shutdown(ctx); err != nil {
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
