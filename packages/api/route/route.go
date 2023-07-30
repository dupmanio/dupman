package route

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewSystemRoute),
		fx.Provide(NewUserRoute),
		fx.Provide(NewWebsiteRoute),
		fx.Invoke(
			func(logger *zap.Logger, systemRoute *SystemRoute, userRoute *UserRoute, websiteRoute *WebsiteRoute) {
				logger.Debug("Setting up routes")

				systemRoute.Setup()
				userRoute.Setup()
				websiteRoute.Setup()
			},
		),
	)
}
