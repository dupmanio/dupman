package notify

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
)

type Deliverer struct {
	notificationSettingsMapping NotificationSettingsMapping
	dupmanCredentials           credentials.Provider
	dupmanAPIService            *service.DupmanAPIService
}

func New(config *config.Config, dupmanAPIService *service.DupmanAPIService) (*Deliverer, error) {
	ctx := context.Background()

	cred, err := credentials.NewClientCredentials(
		ctx,
		config.Dupman.ClientID,
		config.Dupman.ClientSecret,
		[]string{"notify:notification:create"},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman credentials provider: %w", err)
	}

	return &Deliverer{
		notificationSettingsMapping: getNotificationSettingsMapping(),
		dupmanCredentials:           cred,
		dupmanAPIService:            dupmanAPIService,
	}, nil
}

func (del *Deliverer) Name() string {
	return "NotifyDeliverer"
}

func (del *Deliverer) Deliver(message dto.NotificationMessage, _ *dto.ContactInfo) error {
	notificationSettings, ok := del.notificationSettingsMapping[message.Type]
	if !ok {
		return domainErrors.ErrUnsupportedNotificationType
	}

	if err := del.dupmanAPIService.CreateSession(del.dupmanCredentials); err != nil {
		return fmt.Errorf("unable to create dupman session: %w", err)
	}

	_, err := del.dupmanAPIService.NotifySvc.Create(&dto.NotificationOnCreate{
		UserID:  message.UserID,
		Type:    notificationSettings.Type,
		Title:   notificationSettings.Title,
		Message: notificationSettings.Message,
	})
	if err != nil {
		return fmt.Errorf("unable to send Notification: %w", err)
	}

	return nil
}
