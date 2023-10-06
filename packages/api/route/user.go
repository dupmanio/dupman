package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/api/middleware"
	"github.com/dupmanio/dupman/packages/api/server"
	"go.uber.org/zap"
)

type UserRoute struct {
	server         *server.Server
	controller     *controller.UserController
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

func NewUserRoute(
	server *server.Server,
	controller *controller.UserController,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *UserRoute {
	return &UserRoute{
		server:         server,
		controller:     controller,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (route *UserRoute) Setup() {
	route.logger.Debug("Setting up User route")

	group := route.server.Engine.Group(
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
