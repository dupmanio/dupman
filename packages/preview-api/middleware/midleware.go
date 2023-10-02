package middleware

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewAuthMiddleware),
		fx.Provide(NewCORSMiddleware),
		fx.Invoke(
			func(logger *zap.Logger, CORSMiddleware *CORSMiddleware) {
				logger.Debug("Setting up middlewares")

				CORSMiddleware.Setup()
			},
		),
	)
}
