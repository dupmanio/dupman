package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemRoute struct {
	controller  *controller.SystemController
	httpService *commonServices.HTTPService
	logger      *zap.Logger
}

func NewSystemRoute(
	controller *controller.SystemController,
	httpService *commonServices.HTTPService,
	logger *zap.Logger,
) *SystemRoute {
	return &SystemRoute{
		controller:  controller,
		httpService: httpService,
		logger:      logger,
	}
}

func (route *SystemRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up route", zap.String(string(otel.RouteKey), "system"))

	group := engine.Group("/system")
	{
		group.GET("/websites", route.controller.GetWebsites)
		group.POST("/websites/:id/status", route.controller.UpdateWebsiteStatus)
	}
}
