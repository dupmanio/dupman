package repository

import (
	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/model"
	"go.uber.org/zap"
)

type WebsiteRepository struct {
	db     *database.Database
	logger *zap.Logger
}

func NewWebsiteRepository(
	db *database.Database,
	logger *zap.Logger,
) *WebsiteRepository {
	return &WebsiteRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *WebsiteRepository) Setup() {
	repo.logger.Debug("Setting up Website repository")

	if err := repo.db.AutoMigrate(&model.Website{}); err != nil {
		repo.logger.Error(err.Error())
	}
}

func (repo *WebsiteRepository) Create(website *model.Website) {
	// @todo: refactor.
	repo.db.Create(website)
}
