package controller

import "go.uber.org/fx"

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewUserController),
		fx.Provide(NewWebsiteController),
	)
}
