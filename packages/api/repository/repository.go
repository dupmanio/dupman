package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewKeyPairRepository),
		fx.Provide(NewUserRepository),
		fx.Provide(NewWebsiteRepository),
		fx.Invoke(
			func(logger *zap.Logger, keyPairRepo *KeyPairRepository, userRepo *UserRepository, websiteRepo *WebsiteRepository) {
				logger.Debug("Setting up repositories")

				keyPairRepo.Setup()
				userRepo.Setup()
				websiteRepo.Setup()
			},
		),
	)
}
