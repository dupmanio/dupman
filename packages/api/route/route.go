package route

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewWebsiteRoute),
		fx.Invoke(
			func(logger *zap.Logger, websiteRoute *WebsiteRoute) {
				logger.Debug("Setting up routes")

				websiteRoute.Setup()
			},
		),
	)
}
