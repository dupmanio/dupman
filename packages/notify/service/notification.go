package service

import (
	"context"
	"fmt"
	"net"

	authConstant "github.com/dupmanio/dupman/packages/auth/constant"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/notify/config"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/dupmanio/dupman/packages/notify/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	redisClient      *redis.Client
}

func NewNotificationService(
	notificationRepo *repository.NotificationRepository,
	config *config.Config,
) (*NotificationService, error) {
	redisOptions, err := redis.ParseURL(fmt.Sprintf(
		"redis://%s:%s@%s/%s",
		config.Redis.User,
		config.Redis.Password,
		net.JoinHostPort(config.Redis.Host, config.Redis.Port),
		config.Redis.Database,
	))
	if err != nil {
		return nil, fmt.Errorf("unable to estabish redis connection: %w", err)
	}

	return &NotificationService{
		notificationRepo: notificationRepo,
		redisClient:      redis.NewClient(redisOptions),
	}, nil
}

func (svc *NotificationService) Create(ctx *gin.Context, entity *model.Notification) (*model.Notification, error) {
	if err := svc.notificationRepo.Create(ctx.Request.Context(), entity); err != nil {
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

	notifications, err := svc.notificationRepo.FindByUserID(ctx.Request.Context(), userID, pagination)
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

	count, err := svc.notificationRepo.GetCountByUserID(ctx.Request.Context(), userID)
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

	if err := svc.notificationRepo.DeleteByUserID(ctx.Request.Context(), userID); err != nil {
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

	if err := svc.notificationRepo.MarkAsReadByIDAndUserID(ctx.Request.Context(), id, userID); err != nil {
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

	if err := svc.notificationRepo.MarkAsReadByUserID(ctx.Request.Context(), userID); err != nil {
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

func (svc *NotificationService) SendNotificationToChannel(notification dto.NotificationOnResponse) error {
	backgroundContext := context.Background()
	channelName := svc.getUserNotificationsChannelName(notification.UserID)

	if err := svc.redisClient.Publish(backgroundContext, channelName, notification).Err(); err != nil {
		return fmt.Errorf("unable to publish notification: %w", err)
	}

	return nil
}

func (svc *NotificationService) getUserNotificationsChannelName(userID uuid.UUID) string {
	return fmt.Sprintf("user:%s:notifications", userID)
}

func (svc *NotificationService) SubscribeToUserNotifications(ctx *gin.Context) (*redis.PubSub, error) {
	backgroundContext := context.Background()

	userID, err := svc.getCurrentUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get user ID: %w", err)
	}

	channelName := svc.getUserNotificationsChannelName(userID)

	return svc.redisClient.Subscribe(backgroundContext, channelName), nil
}
