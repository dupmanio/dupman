package route

import (
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/notify/controller"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationRoute struct {
	controller  *controller.NotificationController
	httpService *commonServices.HTTPService
	logger      *zap.Logger
}

func NewNotificationRoute(
	controller *controller.NotificationController,
	httpService *commonServices.HTTPService,
	logger *zap.Logger,
) *NotificationRoute {
	return &NotificationRoute{
		controller:  controller,
		httpService: httpService,
		logger:      logger,
	}
}

func (route *NotificationRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up route", zap.String(string(otel.RouteKey), "notification"))

	group := engine.Group("/notification")
	{
		group.POST("", route.controller.Create)
		group.GET("", route.controller.GetAll)
		group.GET("/count", route.controller.GetCount)
		group.GET("/realtime", route.controller.Realtime)
		group.POST("/:id/mark-as-read", route.controller.MarkAsRead)
		group.POST("/mark-all-as-read", route.controller.MarkAllAsRead)
		group.DELETE("", route.controller.DeleteAll)
	}
}
