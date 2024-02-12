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

	authMiddleware := auth.NewMiddleware(
		auth.WithHTTPErrorHandler(route.httpService.HTTPError),
	)
	group := engine.Group("/website")
	{
		group.GET(
			"",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("website-read"),
					filter.NewScopeFilter("api", "api:website", "api:website:read"),
				),
			),
			route.controller.GetAll,
		)
		group.POST(
			"",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("website-create"),
					filter.NewScopeFilter("api", "api:website", "api:website:create"),
				),
			),
			route.controller.Create,
		)
		group.GET(
			"/:id",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("website-read"),
					filter.NewScopeFilter("api", "api:website", "api:website:read"),
				),
			),
			route.controller.GetSingle,
		)
		group.PATCH(
			"",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("website-update"),
					filter.NewScopeFilter("api", "api:website", "api:website:update"),
				),
			),
			route.controller.Update,
		)
		group.DELETE(
			"/:id",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("website-delete"),
					filter.NewScopeFilter("api", "api:website", "api:website:delete"),
				),
			),
			route.controller.Delete,
		)
	}
}
