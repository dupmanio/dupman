package repository

import (
	"context"

	"github.com/dupmanio/dupman/packages/api/constant"
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

func (repo *WebsiteRepository) Create(website *model.Website, encryptionKKey string) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.EncryptionKeyKey, encryptionKKey)

	return repo.db.WithContext(ctx).Create(website).Error
}

func (repo *WebsiteRepository) FindAll() ([]model.Website, error) {
	var websites []model.Website

	return websites, repo.db.Find(&websites).Error
}

func (repo *WebsiteRepository) FindByUserID(userID string) ([]model.Website, error) {
	var websites []model.Website

	return websites, repo.db.Where("user_id", userID).Find(&websites).Error
}
