package middleware

import (
	"github.com/dupmanio/dupman/packages/notify/config"
	"github.com/dupmanio/dupman/packages/notify/server"
	"github.com/gin-contrib/cors"
	"go.uber.org/zap"
)

type CORSMiddleware struct {
	server *server.Server
	logger *zap.Logger
	config *config.Config
}

func NewCORSMiddleware(server *server.Server, logger *zap.Logger, config *config.Config) *CORSMiddleware {
	return &CORSMiddleware{
		server: server,
		logger: logger,
		config: config,
	}
}

func (mid *CORSMiddleware) Setup() {
	mid.logger.Debug("Setting up CORS middleware")

	mid.server.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     mid.config.CORS.AllowOrigins,
		AllowMethods:     mid.config.CORS.AllowMethods,
		AllowHeaders:     mid.config.CORS.AllowHeaders,
		ExposeHeaders:    mid.config.CORS.ExposeHeaders,
		AllowCredentials: true,
	}))
}
