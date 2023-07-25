package service

import (
	"go.uber.org/fx"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(NewHTTPService),
		fx.Provide(NewWebsiteService),
	)
}
