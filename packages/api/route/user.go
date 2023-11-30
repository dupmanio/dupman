package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserRoute struct {
	controller     *controller.UserController
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

func NewUserRoute(
	controller *controller.UserController,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *UserRoute {
	return &UserRoute{
		controller:     controller,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (route *UserRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up User route")

	group := engine.Group(
		"/user",
		route.authMiddleware.RequiresAuth(),
		route.authMiddleware.RequiresRole("service"),
	)
	{
		group.POST(
			"/",
			route.authMiddleware.RequiresScope("user", "user:create"),
			route.controller.Create,
		)
		group.PATCH(
			"/",
			route.authMiddleware.RequiresScope("user", "user:update"),
			route.controller.Update,
		)
		group.GET(
			"/contact-info/:id",
			route.authMiddleware.RequiresScope("user", "user:get_contact_info"),
			route.controller.GetContactInfo,
		)
	}
}
