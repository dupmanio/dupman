package route

import (
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/user-api/controller"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserRoute struct {
	controller  *controller.UserController
	httpService *commonServices.HTTPService
	logger      *zap.Logger
}

func NewUserRoute(
	controller *controller.UserController,
	httpService *commonServices.HTTPService,
	logger *zap.Logger,
) *UserRoute {
	return &UserRoute{
		controller:  controller,
		httpService: httpService,
		logger:      logger,
	}
}

func (route *UserRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up route", zap.String(string(otel.RouteKey), "user"))

	group := engine.Group("/user")
	{
		group.POST("", route.controller.Create)
		group.PATCH("", route.controller.Update)
		group.GET("/contact-info/:id", route.controller.GetContactInfo)
		group.GET("/me", route.controller.Me)
	}
}
