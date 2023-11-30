package service

import (
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		commonServices.NewHTTPService,
		commonServices.NewAuthService,
		NewNotificationService,
	)
}
