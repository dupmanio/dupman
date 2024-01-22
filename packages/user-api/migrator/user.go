package migrator

import (
	"github.com/dupmanio/dupman/packages/common/database"
	"github.com/dupmanio/dupman/packages/user-api/model"
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
