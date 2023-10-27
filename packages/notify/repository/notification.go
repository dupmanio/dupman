package repository

import (
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/notify/database"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationRepository struct {
	db     *database.Database
	logger *zap.Logger
}

func NewNotificationRepository(
	db *database.Database,
	logger *zap.Logger,
) *NotificationRepository {
	return &NotificationRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *NotificationRepository) Setup() {
	repo.logger.Debug("Setting up Notification repository")

	if err := repo.db.AutoMigrate(&model.Notification{}); err != nil {
		repo.logger.Error(err.Error())
	}
}

func (repo *NotificationRepository) Create(notification *model.Notification) error {
	return repo.db.Create(notification).Error
}

func (repo *NotificationRepository) FindByUserID(
	userID uuid.UUID,
	pager *pagination.Pagination,
) ([]model.Notification, error) {
	var notifications []model.Notification

	tx := repo.db.
		Where("user_id", userID).
		Order("created_at DESC")

	return notifications, tx.
		Scopes(pagination.WithPagination(tx, &notifications, pager)).
		Find(&notifications).
		Error
}

func (repo *NotificationRepository) GetCountByUserID(userID uuid.UUID) (int64, error) {
	var count int64

	return count, repo.db.
		Model(&model.Notification{}).
		Where("user_id", userID).
		Where("seen", false).
		Count(&count).
		Error
}

func (repo *NotificationRepository) DeleteByUserID(userID uuid.UUID) error {
	return repo.db.
		Where("user_id", userID).
		Delete(&model.Notification{}).
		Error
}

func (repo *NotificationRepository) MarkAsReadByIDAndUserID(id uuid.UUID, userID uuid.UUID) error {
	return repo.db.
		Model(&model.Notification{}).
		Where("id", id).
		Where("user_id", userID).
		Updates(map[string]interface{}{"Seen": true}).
		Error
}

func (repo *NotificationRepository) MarkAsReadByUserID(userID uuid.UUID) error {
	return repo.db.
		Model(&model.Notification{}).
		Where("user_id", userID).
		Updates(map[string]interface{}{"Seen": true}).
		Error
}
