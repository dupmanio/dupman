package route

import (
	"github.com/dupmanio/dupman/packages/auth"
	"github.com/dupmanio/dupman/packages/auth/filter"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/preview-api/controller"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PreviewRoute struct {
	controller  *controller.PreviewController
	httpService *commonServices.HTTPService
	logger      *zap.Logger
}

func NewPreviewRoute(
	controller *controller.PreviewController,
	httpService *commonServices.HTTPService,
	logger *zap.Logger,
) *PreviewRoute {
	return &PreviewRoute{
		controller:  controller,
		httpService: httpService,
		logger:      logger,
	}
}

func (route *PreviewRoute) Register(engine *gin.Engine) {
	route.logger.Debug("Setting up route", zap.String(string(otel.RouteKey), "preview"))

	authMiddleware := auth.NewMiddleware(
		auth.WithHTTPErrorHandler(route.httpService.HTTPError),
	)

	group := engine.Group("/preview")
	{
		group.GET(
			":id",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("website-get-preview"),
					filter.NewScopeFilter("preview_api", "preview_api:preview", "preview_api:preview:get"),
				),
			),
			route.controller.Preview,
		)
	}
}
