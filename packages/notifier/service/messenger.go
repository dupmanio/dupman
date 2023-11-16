package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/broker"
	"github.com/dupmanio/dupman/packages/notifier/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type MessengerService struct {
	logger *zap.Logger
	config *config.Config
	broker *broker.RabbitMQ
}

func NewMessengerService(
	logger *zap.Logger,
	config *config.Config,
) (*MessengerService, error) {
	brk, err := broker.NewRabbitMQ(&broker.RabbitMQConfig{
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		AppID:    "notifier",
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create RabbitMQ Broker: %w", err)
	}

	if err = brk.Channel.Qos(config.Worker.PrefetchCount, config.Worker.PrefetchSize, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &MessengerService{
		logger: logger,
		config: config,
		broker: brk,
	}, nil
}

func (mess *MessengerService) Consume() (<-chan amqp.Delivery, error) {
	messages, err := mess.broker.ConsumeQueue(mess.config.Worker.QueueName)
	if err != nil {
		return nil, fmt.Errorf("unable to consume messages: %w", err)
	}

	return messages, nil
}

func (mess *MessengerService) AcknowledgeMessage(successfullyProcessed bool, delivery amqp.Delivery) {
	if mess.shouldAcknowledge(successfullyProcessed, delivery) {
		if err := delivery.Ack(false); err != nil {
			mess.logger.Error(
				"Unable to Ack message",
				zap.String("messageID", delivery.MessageId),
				zap.Error(err),
			)
		}

		return
	}

	if err := delivery.Nack(false, false); err != nil {
		mess.logger.Error(
			"Unable to Nack message",
			zap.String("messageID", delivery.MessageId),
			zap.Error(err),
		)
	}
}

func (mess *MessengerService) shouldAcknowledge(successfullyProcessed bool, delivery amqp.Delivery) bool {
	return successfullyProcessed || (!successfullyProcessed && mess.hasReachedMaxRetryCount(delivery))
}

func (mess *MessengerService) hasReachedMaxRetryCount(delivery amqp.Delivery) bool {
	return int(mess.broker.GetDeathCount(delivery)) >= mess.config.Worker.RetryAttempts
}

func (mess *MessengerService) Close() error {
	if err := mess.broker.Close(); err != nil {
		return fmt.Errorf("unable to close RabbitMQ Broker: %w", err)
	}

	return nil
}
