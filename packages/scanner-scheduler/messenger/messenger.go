package messenger

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/broker"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	config *config.Config
	broker *broker.RabbitMQ
	ot     *otel.OTel
}

// @todo: reduce code duplication.
func NewMessengerService(logger *zap.Logger, config *config.Config, ot *otel.OTel) (*Service, error) {
	brk, err := broker.NewRabbitMQ(ot, &broker.RabbitMQConfig{
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		AppID:    config.AppName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create RabbitMQ Broker: %w", err)
	}

	if err = brk.Channel.Confirm(false); err != nil {
		return nil, fmt.Errorf("failed to set Confirm mode: %w", err)
	}

	return &Service{
		logger: logger,
		config: config,
		broker: brk,
		ot:     ot,
	}, nil
}

func (mess *Service) SendScanWebsiteMessage(
	ctx context.Context,
	website dto.WebsiteOnSystemResponse,
) error {
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

func (mess *Service) Close() error {
	if err := mess.broker.Close(); err != nil {
		return fmt.Errorf("unable to close RabbitMQ Broker: %w", err)
	}

	return nil
}
