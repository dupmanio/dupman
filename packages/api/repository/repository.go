package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewWebsiteRepository),
		fx.Invoke(
			func(logger *zap.Logger, websiteRepo *WebsiteRepository) {
				logger.Debug("Setting up repositories")

				websiteRepo.Setup()
			},
		),
	)
}
