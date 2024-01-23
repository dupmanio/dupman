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
		auth.WithCallUserService(false),
		auth.WithHTTPErrorHandler(route.httpService.HTTPError),
		auth.WithFilters(
			filter.NewRoleFilter("user"),
		),
	)
	group := engine.Group("/website")
	{
		group.GET(
			"/",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("website", "website:read"),
				),
			),
			route.controller.GetAll,
		)
		group.POST(
			"/",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("website", "website:create"),
				),
			),
			route.controller.Create,
		)
		group.GET(
			"/:id",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("website", "website:read"),
				),
			),
			route.controller.GetSingle,
		)
		group.PATCH(
			"/",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("website", "website:update"),
				),
			),
			route.controller.Update,
		)
		group.DELETE(
			"/:id",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewScopeFilter("website", "website:delete"),
				),
			),
			route.controller.Delete,
		)
	}
}
