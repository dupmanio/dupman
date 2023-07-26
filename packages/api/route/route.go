package route

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewUserRoute),
		fx.Provide(NewWebsiteRoute),
		fx.Invoke(
			func(logger *zap.Logger, userRoute *UserRoute, websiteRoute *WebsiteRoute) {
				logger.Debug("Setting up routes")

				userRoute.Setup()
				websiteRoute.Setup()
			},
		),
	)
}
