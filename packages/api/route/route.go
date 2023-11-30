package route

import (
	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		fxHelper.AsRoute(NewSystemRoute),
		fxHelper.AsRoute(NewUserRoute),
		fxHelper.AsRoute(NewWebsiteRoute),
	)
}
