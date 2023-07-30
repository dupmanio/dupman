package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewKeyPairRepository),
		fx.Provide(NewUpdateRepository),
		fx.Provide(NewUserRepository),
		fx.Provide(NewWebsiteRepository),
		fx.Invoke(
			func(
				logger *zap.Logger,
				keyPairRepo *KeyPairRepository,
				updateRepo *UpdateRepository,
				userRepo *UserRepository,
				websiteRepo *WebsiteRepository,
			) {
				logger.Debug("Setting up repositories")

				keyPairRepo.Setup()
				updateRepo.Setup()
				userRepo.Setup()
				websiteRepo.Setup()
			},
		),
	)
}
