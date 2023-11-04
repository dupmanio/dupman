package repository

import (
	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/model"
	"go.uber.org/zap"
)

type UpdateRepository struct {
	db     *database.Database
	logger *zap.Logger
}

func NewUpdateRepository(
	db *database.Database,
	logger *zap.Logger,
) *UpdateRepository {
	return &UpdateRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *UpdateRepository) Setup() {
	repo.logger.Debug("Setting up Update repository")

	if err := repo.db.AutoMigrate(&model.Update{}); err != nil {
		repo.logger.Error(err.Error())
	}
}
