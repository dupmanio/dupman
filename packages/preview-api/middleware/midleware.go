package middleware

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewAuthMiddleware),
		fx.Invoke(
			func(logger *zap.Logger) {
				logger.Debug("Setting up middlewares")
			},
		),
	)
}
