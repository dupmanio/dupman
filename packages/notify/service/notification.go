package service

import (
	"context"
	"fmt"
	"net"

	"github.com/dupmanio/dupman/packages/common/ory/keto"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/pagination"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/notify/config"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/dupmanio/dupman/packages/notify/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	redisClient      *redis.Client
	authSvc          *commonServices.AuthService
	ot               *otel.OTel
	keto             *keto.Keto
}

func NewNotificationService(
	notificationRepo *repository.NotificationRepository,
	config *config.Config,
	authSvc *commonServices.AuthService,
	ot *otel.OTel,
	keto *keto.Keto,
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

	redisClient := redis.NewClient(redisOptions)
	if err = ot.InstrumentRedis(redisClient); err != nil {
		return nil, fmt.Errorf("unable to instrumentate redis: %w", err)
	}

	return &NotificationService{
		notificationRepo: notificationRepo,
		redisClient:      redisClient,
		authSvc:          authSvc,
		ot:               ot,
		keto:             keto,
	}, nil
}

func (svc *NotificationService) Create(ctx context.Context, entity *model.Notification) (*model.Notification, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	if err := svc.notificationRepo.Create(ctx, entity); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to create Notification", err)

		return nil, fmt.Errorf("unable to create Notification: %w", err)
	}

	if err := svc.keto.CreateRelationship(
		ctx,
		"Notification",
		entity.ID.String(),
		"recipients",
		entity.UserID.String(),
	); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to Create Notification Recipient Relationship", err)

		return nil, fmt.Errorf("unable to Create Notification Recipient Relationship: %w", err)
	}

	return entity, nil
}

func (svc *NotificationService) GetAllForCurrentUser(
	ctx context.Context,
	pagination *pagination.Pagination,
) ([]model.Notification, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	userID := svc.authSvc.CurrentUserID(ctx)
	if userID == uuid.Nil {
		err := domainErrors.ErrInvalidUserID
		svc.ot.ErrorEvent(ctx, "Invalid User ID", err)

		return nil, err
	}

	notifications, err := svc.notificationRepo.FindByUserID(ctx, userID, pagination)
	if err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to get Notifications", err)

		return nil, fmt.Errorf("unable to get Notifications: %w", err)
	}

	return notifications, nil
}

func (svc *NotificationService) GetCountForCurrentUser(ctx context.Context) (int64, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	userID := svc.authSvc.CurrentUserID(ctx)
	if userID == uuid.Nil {
		err := domainErrors.ErrInvalidUserID
		svc.ot.ErrorEvent(ctx, "Invalid User ID", err)

		return 0, err
	}

	count, err := svc.notificationRepo.GetCountByUserID(ctx, userID)
	if err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to get Notifications count", err)

		return 0, fmt.Errorf("unable to get Notifications count: %w", err)
	}

	return count, nil
}

func (svc *NotificationService) DeleteAllForCurrentUser(ctx context.Context) error {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	userID := svc.authSvc.CurrentUserID(ctx)
	if userID == uuid.Nil {
		err := domainErrors.ErrInvalidUserID
		svc.ot.ErrorEvent(ctx, "Invalid User ID", err)

		return err
	}

	if err := svc.notificationRepo.DeleteByUserID(ctx, userID); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to delete Notifications", err)

		return fmt.Errorf("unable to delete Notifications: %w", err)
	}

	// @todo: list all the affected relationships and remove them.

	return nil
}

func (svc *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	userID := svc.authSvc.CurrentUserID(ctx)
	if userID == uuid.Nil {
		err := domainErrors.ErrInvalidUserID
		svc.ot.ErrorEvent(ctx, "Invalid User ID", err)

		return err
	}

	if err := svc.notificationRepo.MarkAsReadByIDAndUserID(ctx, id, userID); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to mark Notifications as read", err)

		return fmt.Errorf("unable to mark Notifications as read: %w", err)
	}

	return nil
}

func (svc *NotificationService) MarkAllAsReadForCurrentUser(ctx context.Context) error {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	userID := svc.authSvc.CurrentUserID(ctx)
	if userID == uuid.Nil {
		err := domainErrors.ErrInvalidUserID
		svc.ot.ErrorEvent(ctx, "Invalid User ID", err)

		return err
	}

	if err := svc.notificationRepo.MarkAsReadByUserID(ctx, userID); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to mark Notifications as read", err)

		return fmt.Errorf("unable to mark Notifications as read: %w", err)
	}

	return nil
}

func (svc *NotificationService) SendNotificationToChannel(
	ctx context.Context,
	notification dto.NotificationOnResponse,
) error {
	channelName := svc.getUserNotificationsChannelName(notification.UserID)

	if err := svc.redisClient.Publish(ctx, channelName, notification).Err(); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to publish Notification", err)

		return fmt.Errorf("unable to publish Notification: %w", err)
	}

	return nil
}

func (svc *NotificationService) getUserNotificationsChannelName(userID uuid.UUID) string {
	return fmt.Sprintf("user:%s:notifications", userID)
}

func (svc *NotificationService) SubscribeToUserNotifications(ctx context.Context) (*redis.PubSub, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	userID := svc.authSvc.CurrentUserID(ctx)
	if userID == uuid.Nil {
		err := domainErrors.ErrInvalidUserID
		svc.ot.ErrorEvent(ctx, "Invalid User ID", err)

		return nil, err
	}

	channelName := svc.getUserNotificationsChannelName(userID)

	return svc.redisClient.Subscribe(ctx, channelName), nil
}
