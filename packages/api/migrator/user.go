package migrator

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
)

type UserMigrator struct {
	db *database.Database
}

func NewUserMigrator(db *database.Database) *UserMigrator {
	return &UserMigrator{
		db: db,
	}
}

func (mig *UserMigrator) Name() string {
	return "User"
}

func (mig *UserMigrator) Migrate() error {
	return mig.db.AutoMigrate(&model.User{})
}
