package service

import (
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"go.uber.org/fx"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(commonServices.NewDupmanAPIService),
	)
}
