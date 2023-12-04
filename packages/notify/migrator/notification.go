package migrator

import (
	"github.com/dupmanio/dupman/packages/common/database"
	"github.com/dupmanio/dupman/packages/notify/model"
)

type NotificationMigrator struct {
	db *database.Database
}

func NewNotificationMigrator(db *database.Database) *NotificationMigrator {
	return &NotificationMigrator{
		db: db,
	}
}

func (mig *NotificationMigrator) Name() string {
	return "Notification"
}

func (mig *NotificationMigrator) Migrate() error {
	return mig.db.AutoMigrate(&model.Notification{})
}
