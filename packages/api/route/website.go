package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/dupmanio/dupman/packages/api/server"
	"go.uber.org/zap"
)

type WebsiteRoute struct {
	server         *server.Server
	controller     *controller.WebsiteController
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

func NewWebsiteRoute(
	server *server.Server,
	controller *controller.WebsiteController,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *WebsiteRoute {
	return &WebsiteRoute{
		server:         server,
		controller:     controller,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (route *WebsiteRoute) Setup() {
	route.logger.Debug("Setting up Website route")

	group := route.server.Engine.Group("/website")
	{
		group.GET("/", route.authMiddleware.RequiresAuth(), route.controller.GetAll)
		group.POST("/", route.authMiddleware.RequiresAuth(), route.controller.Create)
	}
}
