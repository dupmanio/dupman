package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func (repo *WebsiteRepository) FindAll(pager *pagination.Pagination) ([]model.Website, error) {
	var websites []model.Website

	tx := repo.db.DB

	return websites, tx.
		Scopes(pagination.WithPagination(tx, &websites, pager)).
		Find(&websites).
		Error
}

func (repo *WebsiteRepository) FindByUserID(userID string, pager *pagination.Pagination) ([]model.Website, error) {
	var websites []model.Website

	tx := repo.db.
		Preload("Status").
		Where("websites.user_id", userID)

	return websites, tx.
		Scopes(pagination.WithPagination(tx, &websites, pager)).
		Find(&websites).
		Error
}

func (repo *WebsiteRepository) FindByID(id string) *model.Website {
	var website model.Website

	err := repo.db.
		Preload("Status").
		Preload("Updates").
		First(&website, "id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &website
}

func (repo *WebsiteRepository) ClearUpdates(website *model.Website) error {
	err := repo.db.
		Unscoped().
		Model(&website).
		Association("Updates").
		Unscoped().
		Clear()
	if err != nil {
		return fmt.Errorf("unable to clear Updates Association: %w", err)
	}

	return nil
}

func (repo *WebsiteRepository) UpdateStatus(website *model.Website) error {
	err := repo.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Select("Updates", "Status").
		Omit("UpdatedAt").
		Updates(&website).
		Error
	if err != nil {
		return fmt.Errorf("unable to update Website: %w", err)
	}

	return nil
}

func (repo *WebsiteRepository) DeleteByIDAndUserID(id, userID uuid.UUID) error {
	err := repo.db.
		Debug().
		Unscoped().
		Select("Updates", "Status").
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Delete(&model.Website{
			Base: model.Base{
				ID: id,
			},
		}).
		Error
	if err != nil {
		return fmt.Errorf("unable to delete Website: %w", err)
	}

	return nil
}
