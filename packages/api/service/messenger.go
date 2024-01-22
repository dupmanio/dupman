package service

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/broker"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"go.uber.org/zap"
)

type MessengerService struct {
	logger *zap.Logger
	config *config.Config
	broker *broker.RabbitMQ
	ot     *otel.OTel
}

func NewMessengerService(
	logger *zap.Logger,
	config *config.Config,
	ot *otel.OTel,
) (*MessengerService, error) {
	brk, err := broker.NewRabbitMQ(ot, &broker.RabbitMQConfig{
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		AppID:    "api",
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create RabbitMQ Broker: %w", err)
	}

	if err = brk.Channel.Confirm(false); err != nil {
		return nil, fmt.Errorf("failed to set Confirm mode: %w", err)
	}

	return &MessengerService{
		logger: logger,
		config: config,
		broker: brk,
		ot:     ot,
	}, nil
}

func (mess *MessengerService) SendScanWebsiteMessage(ctx context.Context, website *model.Website) error {
	ctx, span := mess.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	mess.logger.Info(
		"Sending Message",
		zap.String("type", "ScanWebsiteMessage"),
		zap.String("websiteID", website.ID.String()),
	)

	message := dto.ScanWebsiteMessage{
		WebsiteID:    website.ID,
		UserID:       website.UserID,
		WebsiteURL:   website.URL,
		WebsiteToken: website.Token,
	}

	if err := mess.broker.PublishToExchange(
		ctx,
		mess.config.Scanner.ExchangeName,
		mess.config.Scanner.RoutingKey,
		message,
	); err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}

	return nil
}

func (mess *MessengerService) SendStatusChangeNotificationMessage(
	ctx context.Context,
	website *model.Website,
	oldStatus model.Status,
	newStatus model.Status,
	updates []model.Update,
) error {
	ctx, span := mess.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	mess.ot.LogInfoEvent(
		ctx,
		"Sending Message",
		otel.MessageType("StatusChangeNotificationMessage"),
		otel.WebsiteID(website.ID),
	)

	message := dto.NotificationMessage{}

	if newStatus.State == dto.StatusStateNeedsUpdate && oldStatus.State != dto.StatusStateNeedsUpdate {
		message = mess.composeNeedsUpdateNotification(website, updates)
	}

	if newStatus.State == dto.StatusStateScanningFailed && oldStatus.State != dto.StatusStateScanningFailed {
		message = mess.composeScanningFailedNotification(website, newStatus)
	}

	if message.Type != "" {
		if err := mess.broker.PublishToExchange(
			ctx,
			mess.config.Notifier.ExchangeName,
			mess.config.Notifier.RoutingKey,
			message,
		); err != nil {
			mess.ot.ErrorEvent(ctx, "Failed to publish message to exchange", err)

			return fmt.Errorf("failed to publish message to exchange: %w", err)
		}
	}

	return nil
}

func (mess *MessengerService) composeNeedsUpdateNotification(
	website *model.Website,
	updates []model.Update,
) dto.NotificationMessage {
	updatesMapping := map[string]string{}
	for _, update := range updates {
		updatesMapping[update.Name] = update.RecommendedVersion
	}

	return dto.NotificationMessage{
		UserID: website.UserID,
		Type:   "WebsiteNeedsUpdates",
		Meta: map[string]any{
			"WebsiteID":  website.ID,
			"WebsiteURL": website.URL,
			"Updates":    updatesMapping,
		},
	}
}

func (mess *MessengerService) composeScanningFailedNotification(
	website *model.Website,
	status model.Status,
) dto.NotificationMessage {
	return dto.NotificationMessage{
		UserID: website.UserID,
		Type:   "WebsiteScanningFailed",
		Meta: map[string]any{
			"WebsiteID":  website.ID,
			"WebsiteURL": website.URL,
			"StatusInfo": status.Info,
		},
	}
}

func (mess *MessengerService) Close() error {
	if err := mess.broker.Close(); err != nil {
		return fmt.Errorf("unable to close RabbitMQ Broker: %w", err)
	}

	return nil
}
