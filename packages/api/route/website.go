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

	group := route.server.Engine.Group(
		"/website",
		route.authMiddleware.RequiresAuth(),
		route.authMiddleware.RequiresRole("user"),
	)
	{
		group.GET(
			"/",
			route.authMiddleware.RequiresScope("website", "website:read"),
			route.controller.GetAll,
		)
		group.POST(
			"/",
			route.authMiddleware.RequiresScope("website", "website:create"),
			route.controller.Create,
		)
		group.GET(
			"/:id",
			route.authMiddleware.RequiresScope("website", "website:read"),
			route.controller.GetSingle,
		)
		group.PATCH(
			"/",
			route.authMiddleware.RequiresScope("website", "website:update"),
			route.controller.Update,
		)
		group.DELETE(
			"/:id",
			route.authMiddleware.RequiresScope("website", "website:delete"),
			route.controller.Delete,
		)
	}
}
