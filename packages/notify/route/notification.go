package route

import (
	"github.com/dupmanio/dupman/packages/auth"
	"github.com/dupmanio/dupman/packages/auth/filter"
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

	authMiddleware := auth.NewMiddleware(
		auth.WithCallUserService(false),
		auth.WithHTTPErrorHandler(route.httpService.HTTPError),
	)

	group := engine.Group("/notification")
	{
		group.POST(
			"/",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("service"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:create"),
				),
			),
			route.controller.Create,
		)
		group.GET(
			"/",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:read"),
				),
			),
			route.controller.GetAll,
		)
		group.GET(
			"/count",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:read"),
				),
			),
			route.controller.GetCount,
		)
		group.GET(
			"/realtime",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:read"),
				),
			),
			route.controller.Realtime,
		)
		group.POST(
			"/:id/mark-as-read",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:update"),
				),
			),
			route.controller.MarkAsRead,
		)
		group.POST(
			"/mark-all-as-read",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:update"),
				),
			),
			route.controller.MarkAllAsRead,
		)
		group.DELETE(
			"/",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user"),
					filter.NewScopeFilter("notify", "notify:notification", "notify:notification:delete"),
				),
			),
			route.controller.DeleteAll,
		)
	}
}
