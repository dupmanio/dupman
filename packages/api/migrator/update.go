package migrator

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
)

type UpdateMigrator struct {
	db *database.Database
}

func NewUpdateMigrator(db *database.Database) *UpdateMigrator {
	return &UpdateMigrator{
		db: db,
	}
}

func (mig *UpdateMigrator) Name() string {
	return "Update"
}

func (mig *UpdateMigrator) Migrate() error {
	return mig.db.AutoMigrate(&model.Update{})
}
