package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebsiteRoute struct {
	controller     *controller.WebsiteController
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

func NewWebsiteRoute(
	controller *controller.WebsiteController,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *WebsiteRoute {
	return &WebsiteRoute{
		controller:     controller,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (route *WebsiteRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up Website route")

	group := engine.Group(
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
