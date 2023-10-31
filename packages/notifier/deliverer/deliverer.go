package deliverer

import (
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/email"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/notify"
	"go.uber.org/fx"
)

type Deliverer interface {
	Name() string
	Deliver(dto.NotificationMessage, *dto.ContactInfo) error
}

func Create() fx.Option {
	return fx.Options(
		fx.Provide(email.New),
		fx.Provide(notify.New),
	)
}
