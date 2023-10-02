package route

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewPreviewRoute),
		fx.Invoke(
			func(logger *zap.Logger, userRoute *PreviewRoute) {
				logger.Debug("Setting up routes")

				userRoute.Setup()
			},
		),
	)
}
