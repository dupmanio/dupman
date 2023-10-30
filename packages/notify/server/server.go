package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dupmanio/dupman/packages/notify/config"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	Engine    *gin.Engine
	Interrupt chan os.Signal
}

func New(logger *zap.Logger, config *config.Config) (*Server, error) {
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

	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))

	return &Server{
		Engine:    engine,
		Interrupt: make(chan os.Signal, 1),
	}, nil
}

func Run(server *Server, lc fx.Lifecycle, logger *zap.Logger, config *config.Config) error {
	httpServer := http.Server{
		Addr:              net.JoinHostPort(config.Server.ListenAddr, config.Server.Port),
		Handler:           server.Engine,
		ReadHeaderTimeout: 0,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP Server", zap.String("address", httpServer.Addr))

			go func() {
				if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP Server")

			server.Interrupt <- syscall.SIGTERM
			if err := httpServer.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown server: %w", err)
			}

			return nil
		},
	})

	return nil
}
