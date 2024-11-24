package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebsiteRoute struct {
	controller  *controller.WebsiteController
	httpService *commonServices.HTTPService
	logger      *zap.Logger
}

func NewWebsiteRoute(
	controller *controller.WebsiteController,
	httpService *commonServices.HTTPService,
	logger *zap.Logger,
) *WebsiteRoute {
	return &WebsiteRoute{
		controller:  controller,
		httpService: httpService,
		logger:      logger,
	}
}

func (route *WebsiteRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up route", zap.String(string(otel.RouteKey), "website"))

	group := engine.Group("/website")
	{
		group.GET("", route.controller.GetAll)
		group.POST("", route.controller.Create)
		group.GET("/:id", route.controller.GetSingle)
		group.PATCH("/:id", route.controller.Update)
		group.DELETE("/:id", route.controller.Delete)
	}
}
