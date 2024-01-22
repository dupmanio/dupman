package repository

import (
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		NewWebsiteRepository,
	)
}
