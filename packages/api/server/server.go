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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	Engine *gin.Engine
}

func New(logger *zap.Logger, config *config.Config) (*Server, error) {
	ginLogWrapper := logWrapper.NewGinWrapper(logger)

	gin.DefaultWriter = ginLogWrapper

	if config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.Env == "test" {
		gin.SetMode(gin.TestMode)
	}

	engine := gin.New()

	if err := engine.SetTrustedProxies(config.Server.TrustedProxies); err != nil {
		return nil, fmt.Errorf("unable to set trusted proxies: %w", err)
	}

	if config.Telemetry.Enabled {
		engine.Use(otelgin.Middleware(config.AppName))
	}

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
) error {
	httpServer := http.Server{
		Addr:              net.JoinHostPort(config.Server.ListenAddr, config.Server.Port),
		Handler:           server.Engine,
		ReadHeaderTimeout: 0,
	}

	var ot *otel.OTel

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error

			logger.Info("Starting HTTP Server", zap.String("address", httpServer.Addr))

			if config.Telemetry.Enabled {
				ot, err = otel.NewOTel(ctx, config.AppName, "1.0.0", config.Telemetry.CollectorURL)
				if err != nil {
					return fmt.Errorf("failed to initialize Telemetry service: %w", err)
				}
			}

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

			if config.Telemetry.Enabled {
				if err := ot.Shutdown(ctx); err != nil {
					return fmt.Errorf("failed to shutdown telemetry service: %w", err)
				}
			}

			return nil
		},
	})

	return nil
}
