package messenger

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/broker"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	config *config.Config
	broker *broker.RabbitMQ
}

func NewMessengerService(logger *zap.Logger, config *config.Config) (*Service, error) {
	brk, err := broker.NewRabbitMQ(&broker.RabbitMQConfig{
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		AppID:    "scanner-scheduler",
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
	}, nil
}

func (mess *Service) SendScanWebsiteMessage(website dto.WebsiteOnSystemResponse, token string) error {
	mess.logger.Info(
		"Sending Message",
		zap.String("type", "ScanWebsiteMessage"),
		zap.String("websiteID", website.ID.String()),
	)

	message := dto.ScanWebsiteMessage{
		WebsiteID:    website.ID,
		WebsiteURL:   website.URL,
		WebsiteToken: token,
	}

	if err := mess.broker.PublishToExchange(
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
