package route

import (
	"github.com/dupmanio/dupman/packages/auth"
	"github.com/dupmanio/dupman/packages/auth/filter"
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

	authMiddleware := auth.NewMiddleware(
		auth.WithHTTPErrorHandler(route.httpService.HTTPError),
	)

	group := engine.Group("/user")
	{
		group.POST(
			"",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user-create"),
					filter.NewScopeFilter("user_api", "user_api:user", "user_api:user:create"),
				),
			),
			route.controller.Create,
		)
		group.PATCH(
			"",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user-update"),
					filter.NewScopeFilter("user_api", "user_api:user", "user_api:user:update"),
				),
			),
			route.controller.Update,
		)
		group.GET(
			"/contact-info/:id",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user-get-contact-info"),
					filter.NewScopeFilter("user_api", "user_api:user", "user_api:user:get_contact_info"),
				),
			),
			route.controller.GetContactInfo,
		)
		group.GET(
			"/me",
			authMiddleware.Handler(
				auth.WithFilters(
					filter.NewRoleFilter("user-read"),
					filter.NewScopeFilter("user_api", "user_api:user", "user_api:user:me"),
				),
			),
			route.controller.Me,
		)
	}
}
