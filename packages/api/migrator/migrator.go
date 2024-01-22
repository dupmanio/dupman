package migrator

import (
	fxHelper "github.com/dupmanio/dupman/packages/common/helper/fx"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		fxHelper.AsMigrator(NewStatusMigrator),
		fxHelper.AsMigrator(NewUpdateMigrator),
		fxHelper.AsMigrator(NewWebsiteMigrator),
	)
}
