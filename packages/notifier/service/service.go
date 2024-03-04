package service

import (
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		NewMessengerService,
	)
}
