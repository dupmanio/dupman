package fx

import (
	"github.com/dupmanio/dupman/packages/common/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type IMigrator interface {
	Name() string
	Migrate() error
}

const migratorGroupTag = `group:"migrators"`

func AsMigrator(fc any, annotations ...fx.Annotation) any {
	annotations = append(
		annotations,
		fx.As(new(IMigrator)),
		fx.ResultTags(migratorGroupTag),
	)

	return fx.Annotate(fc, annotations...)
}

func Migrate(logger *zap.Logger, migrators ...IMigrator) {
	for _, migrator := range migrators {
		logger.Debug("Executing Migration", zap.String(string(otel.MigratorKey), migrator.Name()))

		if err := migrator.Migrate(); err != nil {
			logger.Error(
				"Unable to execute migration",
				zap.String(string(otel.MigratorKey), migrator.Name()),
				zap.Error(err),
			)
		}
	}
}
