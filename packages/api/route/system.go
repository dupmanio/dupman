package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/dupmanio/dupman/packages/api/server"
	"go.uber.org/zap"
)

type SystemRoute struct {
	server         *server.Server
	controller     *controller.SystemController
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

func NewSystemRoute(
	server *server.Server,
	controller *controller.SystemController,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *SystemRoute {
	return &SystemRoute{
		server:         server,
		controller:     controller,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (route *SystemRoute) Setup() {
	route.logger.Debug("Setting up System route")

	group := route.server.Engine.Group(
		"/system",
		route.authMiddleware.RequiresAuth(),
		route.authMiddleware.RequiresRole("service"),
	)
	{
		group.GET(
			"/websites",
			route.authMiddleware.RequiresScope("system", "system:get_websites"),
			route.controller.GetWebsites,
		)
	}
}
