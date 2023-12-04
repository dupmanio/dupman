package repository

import (
	"context"

	"github.com/dupmanio/dupman/packages/common/database"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/google/uuid"
)

type NotificationRepository struct {
	db *database.Database
}

func NewNotificationRepository(db *database.Database) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (repo *NotificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	return repo.db.
		WithContext(ctx).
		Create(notification).
		Error
}

func (repo *NotificationRepository) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	pager *pagination.Pagination,
) ([]model.Notification, error) {
	var notifications []model.Notification

	tx := repo.db.
		WithContext(ctx).
		Where("user_id", userID).
		Order("created_at DESC")

	return notifications, tx.
		Scopes(pagination.WithPagination(tx, &notifications, pager)).
		Find(&notifications).
		Error
}

func (repo *NotificationRepository) GetCountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64

	return count, repo.db.
		WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id", userID).
		Where("seen", false).
		Count(&count).
		Error
}

func (repo *NotificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return repo.db.
		WithContext(ctx).
		Where("user_id", userID).
		Delete(&model.Notification{}).
		Error
}

func (repo *NotificationRepository) MarkAsReadByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return repo.db.
		WithContext(ctx).
		Model(&model.Notification{}).
		Where("id", id).
		Where("user_id", userID).
		Updates(map[string]interface{}{"Seen": true}).
		Error
}

func (repo *NotificationRepository) MarkAsReadByUserID(ctx context.Context, userID uuid.UUID) error {
	return repo.db.
		WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id", userID).
		Updates(map[string]interface{}{"Seen": true}).
		Error
}
