package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
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

func (repo *WebsiteRepository) Create(ctx context.Context, website *model.Website, encryptionKKey string) error {
	ctx = context.WithValue(ctx, constant.EncryptionKeyKey, encryptionKKey)

	return repo.db.
		WithContext(ctx).
		Create(website).
		Error
}

func (repo *WebsiteRepository) FindAll(ctx context.Context, pager *pagination.Pagination) ([]model.Website, error) {
	var websites []model.Website

	tx := repo.db.DB

	return websites, tx.
		WithContext(ctx).
		Scopes(pagination.WithPagination(tx, &websites, pager)).
		Find(&websites).
		Error
}

func (repo *WebsiteRepository) FindByUserID(
	ctx context.Context,
	userID string,
	pager *pagination.Pagination,
) ([]model.Website, error) {
	var websites []model.Website

	tx := repo.db.
		WithContext(ctx).
		Preload("Status").
		Where("websites.user_id", userID)

	return websites, tx.
		Scopes(pagination.WithPagination(tx, &websites, pager)).
		Find(&websites).
		Error
}

func (repo *WebsiteRepository) FindByID(ctx context.Context, id string) *model.Website {
	var website model.Website

	err := repo.db.
		WithContext(ctx).
		Preload("Status").
		Preload("Updates").
		First(&website, "id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &website
}

func (repo *WebsiteRepository) ClearUpdates(ctx context.Context, website *model.Website) error {
	err := repo.db.
		WithContext(ctx).
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

func (repo *WebsiteRepository) UpdateStatus(ctx context.Context, website *model.Website) error {
	err := repo.db.
		WithContext(ctx).
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

func (repo *WebsiteRepository) Update(
	ctx context.Context,
	website *model.Website,
	fieldsToUpdate []string,
	encryptionKKey string,
) error {
	ctx = context.WithValue(ctx, constant.EncryptionKeyKey, encryptionKKey)

	return repo.db.
		WithContext(ctx).
		Select("UpdatedAt", fieldsToUpdate).
		Save(website).
		Error
}

func (repo *WebsiteRepository) DeleteByIDAndUserID(ctx context.Context, id, userID uuid.UUID) error {
	err := repo.db.
		WithContext(ctx).
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
