package service

import (
	"fmt"

	authConstant "github.com/dupmanio/dupman/packages/auth/constant"
	"github.com/dupmanio/dupman/packages/common/pagination"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/dupmanio/dupman/packages/notify/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(
	notificationRepo *repository.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (svc *NotificationService) Create(entity *model.Notification) (*model.Notification, error) {
	if err := svc.notificationRepo.Create(entity); err != nil {
		return nil, fmt.Errorf("unable to create notification: %w", err)
	}

	return entity, nil
}

func (svc *NotificationService) GetAllForCurrentUser(
	ctx *gin.Context,
	pagination *pagination.Pagination,
) ([]model.Notification, error) {
	userID, err := svc.getCurrentUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get user ID: %w", err)
	}

	notifications, err := svc.notificationRepo.FindByUserID(userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("unable to get notifications: %w", err)
	}

	return notifications, nil
}

func (svc *NotificationService) GetCountForCurrentUser(ctx *gin.Context) (int64, error) {
	userID, err := svc.getCurrentUserID(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to get user ID: %w", err)
	}

	count, err := svc.notificationRepo.GetCountByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("unable to get notifications count: %w", err)
	}

	return count, nil
}

func (svc *NotificationService) DeleteAllForCurrentUser(
	ctx *gin.Context,
) error {
	userID, err := svc.getCurrentUserID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get user ID: %w", err)
	}

	if err := svc.notificationRepo.DeleteByUserID(userID); err != nil {
		return fmt.Errorf("unable to delete notifications: %w", err)
	}

	return nil
}

func (svc *NotificationService) MarkAsRead(
	ctx *gin.Context,
	id uuid.UUID,
) error {
	userID, err := svc.getCurrentUserID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get user ID: %w", err)
	}

	if err := svc.notificationRepo.MarkAsReadByIDAndUserID(id, userID); err != nil {
		return fmt.Errorf("unable to mark notifications as read: %w", err)
	}

	return nil
}

func (svc *NotificationService) MarkAllAsReadForCurrentUser(
	ctx *gin.Context,
) error {
	userID, err := svc.getCurrentUserID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get user ID: %w", err)
	}

	if err := svc.notificationRepo.MarkAsReadByUserID(userID); err != nil {
		return fmt.Errorf("unable to mark notifications as read: %w", err)
	}

	return nil
}

func (svc *NotificationService) getCurrentUserID(ctx *gin.Context) (uuid.UUID, error) {
	if userIDInterface, ok := ctx.Get(authConstant.UserIDKey); ok {
		if userIDRaw, ok := userIDInterface.(string); ok {
			userID, err := uuid.Parse(userIDRaw)
			if err != nil {
				return uuid.Nil, fmt.Errorf("unable to parse ID: %w", err)
			}

			return userID, nil
		}
	}

	return uuid.Nil, domainErrors.ErrInvalidUserID
}
