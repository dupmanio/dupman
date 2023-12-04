package migrator

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
)

type StatusMigrator struct {
	db *database.Database
}

func NewStatusMigrator(db *database.Database) *StatusMigrator {
	return &StatusMigrator{
		db: db,
	}
}

func (mig *StatusMigrator) Name() string {
	return "Status"
}

func (mig *StatusMigrator) Migrate() error {
	mig.db.Exec(`
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_state') THEN
		CREATE TYPE status_state AS ENUM('UP_TO_DATED', 'NEEDS_UPDATE', 'SCANNING_FAILED');
    END IF;
END $$;
`)

	return mig.db.AutoMigrate(&model.Status{})
}
