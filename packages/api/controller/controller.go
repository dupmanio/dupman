package controller

import "go.uber.org/fx"

func Create() fx.Option {
	return fx.Provide(
		NewSystemController,
		NewUserController,
		NewWebsiteController,
	)
}
