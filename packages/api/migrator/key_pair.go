package migrator

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
)

type KeyPairMigrator struct {
	db *database.Database
}

func NewKeyPairMigrator(db *database.Database) *KeyPairMigrator {
	return &KeyPairMigrator{
		db: db,
	}
}

func (mig *KeyPairMigrator) Name() string {
	return "KeyPair"
}

func (mig *KeyPairMigrator) Migrate() error {
	return mig.db.AutoMigrate(&model.KeyPair{})
}
