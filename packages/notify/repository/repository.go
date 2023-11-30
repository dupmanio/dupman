package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewNotificationRepository),
		fx.Invoke(
			func(
				logger *zap.Logger,
				notificationRepo *NotificationRepository,
			) {
				logger.Debug("Setting up repositories")

				notificationRepo.Setup()
			},
		),
	)
}
