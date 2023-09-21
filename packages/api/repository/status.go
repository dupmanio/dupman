package repository

import (
	"errors"

	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_state') THEN
        DROP TYPE status_state CASCADE;
    END IF;

	CREATE TYPE status_state AS ENUM('UP_TO_DATED', 'NEEDS_UPDATE', 'SCANNING_FAILED');
END $$;
`)

	if err := repo.db.AutoMigrate(&model.Status{}); err != nil {
		repo.logger.Error(err.Error())
	}
}

func (repo *StatusRepository) Create(status *model.Status) error {
	return repo.db.Create(status).Error
}

func (repo *StatusRepository) Update(status *model.Status) error {
	return repo.db.Save(status).Error
}

func (repo *StatusRepository) FindByWebsiteID(id string) *model.Status {
	var status model.Status

	err := repo.db.First(&status, "website_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &status
}
