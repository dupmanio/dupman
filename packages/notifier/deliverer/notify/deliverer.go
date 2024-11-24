package notify

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/service/notify"
)

type Deliverer struct {
	notificationSettingsMapping NotificationSettingsMapping
	dupmanCredentials           credentials.Provider
	config                      *config.Config
}

func New(config *config.Config) (*Deliverer, error) {
	ctx := context.Background()

	cred, err := credentials.NewClientCredentials(
		ctx,
		config.Dupman.ClientID,
		config.Dupman.ClientSecret,
		config.Dupman.Scopes,
		config.Dupman.Audience,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman credentials provider: %w", err)
	}

	return &Deliverer{
		notificationSettingsMapping: getNotificationSettingsMapping(),
		dupmanCredentials:           cred,
		config:                      config,
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

	notifySvc := notify.New(dupman.NewConfig(
		dupman.WithBaseURL(del.config.ServiceURL.Notify),
		dupman.WithCredentials(del.dupmanCredentials),
		dupman.WithOTelEnabled(),
	))

	_, err := notifySvc.Create(&dto.NotificationOnCreate{
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
