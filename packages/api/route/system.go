package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemRoute struct {
	controller     *controller.SystemController
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

func NewSystemRoute(
	controller *controller.SystemController,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *SystemRoute {
	return &SystemRoute{
		controller:     controller,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (route *SystemRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up System route")

	group := engine.Group(
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
		group.POST(
			"/websites/:id/status",
			route.authMiddleware.RequiresScope("system", "system:update_website_status"),
			route.controller.UpdateWebsiteStatus,
		)
	}
}
