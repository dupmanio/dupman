package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/config"
	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(
	env string,
	config config.ServerConfig,
	routes []fxHelper.IRoute,
	logger *zap.Logger,
	ot *otel.OTel,
) (*http.Server, error) {
	ginLogWrapper := logWrapper.NewGinWrapper(logger)

	gin.DefaultWriter = ginLogWrapper

	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	if env == "test" {
		gin.SetMode(gin.TestMode)
	}

	engine := gin.New()
	engine.ContextWithFallback = true

	if err := engine.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, fmt.Errorf("unable to set trusted proxies: %w", err)
	}

	engine.Use(ot.GetOTelGinMiddleware())
	engine.Use(ginLogWrapper.GetGinzapMiddleware())
	engine.Use(ginLogWrapper.GetGinzapRecoveryMiddleware())

	fxHelper.RegisterRoutes(engine, routes...)

	return &http.Server{
		Addr:              net.JoinHostPort(config.ListenAddr, config.Port),
		Handler:           engine,
		ReadHeaderTimeout: 0,
	}, nil
}
