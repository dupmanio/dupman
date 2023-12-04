package migrator

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
)

type WebsiteMigrator struct {
	db *database.Database
}

func NewWebsiteMigrator(db *database.Database) *WebsiteMigrator {
	return &WebsiteMigrator{
		db: db,
	}
}

func (mig *WebsiteMigrator) Name() string {
	return "Website"
}

func (mig *WebsiteMigrator) Migrate() error {
	return mig.db.AutoMigrate(&model.Website{})
}
