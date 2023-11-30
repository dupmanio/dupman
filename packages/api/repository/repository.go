package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewKeyPairRepository),
		fx.Provide(NewStatusRepository),
		fx.Provide(NewUpdateRepository),
		fx.Provide(NewUserRepository),
		fx.Provide(NewWebsiteRepository),
		fx.Invoke(
			func(
				logger *zap.Logger,
				keyPairRepo *KeyPairRepository,
				statusRepo *StatusRepository,
				updateRepo *UpdateRepository,
				userRepo *UserRepository,
				websiteRepo *WebsiteRepository,
			) {
				logger.Debug("Setting up repositories")

				keyPairRepo.Setup()
				updateRepo.Setup()
				statusRepo.Setup()
				userRepo.Setup()
				websiteRepo.Setup()
			},
		),
	)
}
