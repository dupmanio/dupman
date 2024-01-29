package deliverer

import (
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/email"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/notify"
	"go.uber.org/fx"
)

type Deliverer interface {
	Name() string
	Deliver(message dto.NotificationMessage, contactInfo *dto.ContactInfo) error
}

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotate(
			email.New,
			fx.As(new(Deliverer)),
			fx.ResultTags(`group:"deliverers"`),
		),
		fx.Annotate(
			notify.New,
			fx.As(new(Deliverer)),
			fx.ResultTags(`group:"deliverers"`),
		),
	)
}
