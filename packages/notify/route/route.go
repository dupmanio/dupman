package route

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewNotificationRoute),
		fx.Invoke(
			func(logger *zap.Logger, notificationRoute *NotificationRoute) {
				logger.Debug("Setting up routes")

				notificationRoute.Setup()
			},
		),
	)
}
