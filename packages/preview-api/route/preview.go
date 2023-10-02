package route

import (
	"github.com/dupmanio/dupman/packages/preview-api/controller"
	"github.com/dupmanio/dupman/packages/preview-api/middleware"
	"github.com/dupmanio/dupman/packages/preview-api/server"
	"go.uber.org/zap"
)

type PreviewRoute struct {
	server     *server.Server
	controller *controller.PreviewController
	authMid    *middleware.AuthMiddleware
	logger     *zap.Logger
}

func NewPreviewRoute(
	server *server.Server,
	controller *controller.PreviewController,
	authMid *middleware.AuthMiddleware,
	logger *zap.Logger,
) *PreviewRoute {
	return &PreviewRoute{
		server:     server,
		controller: controller,
		authMid:    authMid,
		logger:     logger,
	}
}

func (route *PreviewRoute) Setup() {
	route.logger.Debug("Setting up Preview route")

	group := route.server.Engine.Group(
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
