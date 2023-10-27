package route

import (
	"github.com/dupmanio/dupman/packages/notify/controller"
	"github.com/dupmanio/dupman/packages/notify/middleware"
	"github.com/dupmanio/dupman/packages/notify/server"
	"go.uber.org/zap"
)

type NotificationRoute struct {
	server     *server.Server
	controller *controller.NotificationController
	authMid    *middleware.AuthMiddleware
	logger     *zap.Logger
}

func NewNotificationRoute(
	server *server.Server,
	controller *controller.NotificationController,
	authMid *middleware.AuthMiddleware,
	logger *zap.Logger,
) *NotificationRoute {
	return &NotificationRoute{
		server:     server,
		controller: controller,
		authMid:    authMid,
		logger:     logger,
	}
}

func (route *NotificationRoute) Setup() {
	route.logger.Debug("Setting up Notification route")

	group := route.server.Engine.Group(
		"/notification",
		route.authMid.RequiresAuth(),
	)
	{
		group.POST(
			"/",
			route.authMid.RequiresRole("service"),
			route.authMid.RequiresScope("notify", "notify:notification", "notify:notification:create"),
			route.controller.Create,
		)
		group.GET(
			"/",
			route.authMid.RequiresRole("user"),
			route.authMid.RequiresScope("notify", "notify:notification", "notify:notification:read"),
			route.controller.GetAll,
		)
		group.GET(
			"/count",
			route.authMid.RequiresRole("user"),
			route.authMid.RequiresScope("notify", "notify:notification", "notify:notification:read"),
			route.controller.GetCount,
		)
		group.POST(
			"/:id/mark-as-read",
			route.authMid.RequiresRole("user"),
			route.authMid.RequiresScope("notify", "notify:notification", "notify:notification:update"),
			route.controller.MarkAsRead,
		)
		group.POST(
			"/mark-all-as-read",
			route.authMid.RequiresRole("user"),
			route.authMid.RequiresScope("notify", "notify:notification", "notify:notification:update"),
			route.controller.MarkAllAsRead,
		)
		group.DELETE(
			"/",
			route.authMid.RequiresRole("user"),
			route.authMid.RequiresScope("notify", "notify:notification", "notify:notification:delete"),
			route.controller.DeleteAll,
		)
	}
}
