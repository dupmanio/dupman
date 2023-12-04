package repository

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
	"go.uber.org/zap"
)

type StatusRepository struct {
	db     *database.Database
	logger *zap.Logger
}

func NewStatusRepository(
	db *database.Database,
	logger *zap.Logger,
) *StatusRepository {
	return &StatusRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *StatusRepository) Setup() {
	repo.logger.Debug("Setting up Status repository")

	repo.db.Exec(`
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_state') THEN
		CREATE TYPE status_state AS ENUM('UP_TO_DATED', 'NEEDS_UPDATE', 'SCANNING_FAILED');
    END IF;
END $$;
`)

	if err := repo.db.AutoMigrate(&model.Status{}); err != nil {
		repo.logger.Error(err.Error())
	}
}
