package route

import (
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/preview-api/controller"
	"github.com/dupmanio/dupman/packages/preview-api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PreviewRoute struct {
	controller *controller.PreviewController
	authMid    *middleware.AuthMiddleware
	logger     *zap.Logger
}

func NewPreviewRoute(
	controller *controller.PreviewController,
	authMid *middleware.AuthMiddleware,
	logger *zap.Logger,
) *PreviewRoute {
	return &PreviewRoute{
		controller: controller,
		authMid:    authMid,
		logger:     logger,
	}
}

func (route *PreviewRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up route", zap.String(string(otel.RouteKey), "preview"))

	group := engine.Group(
		"/preview",
		route.authMid.RequiresAuth(),
		route.authMid.RequiresRole("user"),
	)
	{
		group.GET(
			":id",
			route.authMid.RequiresScope("preview_api", "preview_api:preview", "preview_api:preview:get"),
			route.controller.Preview,
		)
	}
}
