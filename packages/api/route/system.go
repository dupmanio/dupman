package route

import (
	"github.com/dupmanio/dupman/packages/api/controller"
	"github.com/dupmanio/dupman/packages/auth"
	"github.com/dupmanio/dupman/packages/auth/filter"
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

	authMiddleware := auth.NewMiddleware(
		auth.WithHTTPErrorHandler(route.httpService.HTTPError),
		auth.WithFilters(
			filter.NewRoleFilter("service"),
		),
	)

	group := engine.Group("/system")
	{
		group.GET(
			"/websites",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("system", "system:get_websites"),
				),
			),
			route.controller.GetWebsites,
		)
		group.POST(
			"/websites/:id/status",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("system", "system:update_website_status"),
				),
			),
			route.controller.UpdateWebsiteStatus,
		)
	}
}
