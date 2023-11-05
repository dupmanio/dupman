package broker

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	config     *config.Config
	logger     *zap.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(config *config.Config, logger *zap.Logger) (*RabbitMQ, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		net.JoinHostPort(config.RabbitMQ.Host, config.RabbitMQ.Port),
	)

	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	if err := channel.Confirm(false); err != nil {
		return nil, fmt.Errorf("failed to set Confirm mode: %w", err)
	}

	return &RabbitMQ{
		logger:     logger,
		config:     config,
		connection: connection,
		channel:    channel,
	}, nil
}

func (brk *RabbitMQ) PublishToScanner(notification []byte) error {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// @todo: Implement ack checking and retry functionality.
	// pubAck, pubNack := brk.channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))
	// select {
	// case <-pubAck:
	//	 return nil
	// case <-pubNack:
	//	 return domainErrors.ErrUnableToPublishMessage
	// case <-time.After(5 * time.Second):
	//	 return domainErrors.ErrUnableToPublishMessage
	// }
	// @todo: add some logging.
	err := brk.channel.PublishWithContext(ctx,
		brk.config.Scanner.ExchangeName,
		brk.config.Scanner.RoutingKey,
		false,
		false,
		amqp.Publishing{
			MessageId:    uuid.New().String(),
			UserId:       brk.config.RabbitMQ.User,
			AppId:        "scanner-scheduler",
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         notification,
		})
	if err != nil {
		return fmt.Errorf("failed to publish Notification: %w", err)
	}

	return nil
}

func (brk *RabbitMQ) Shutdown() error {
	brk.logger.Info("Closing channel")

	if err := brk.channel.Close(); err != nil {
		brk.logger.Error("Unable to close channel", zap.Error(err))

		return fmt.Errorf("unable to close channel: %w", err)
	}

	brk.logger.Info("Channel has been Closed")
	brk.logger.Info("Closing connection")

	if err := brk.connection.Close(); err != nil {
		brk.logger.Error("Unable to close connection", zap.Error(err))

		return fmt.Errorf("unable to close connection: %w", err)
	}

	brk.logger.Info("Connection has been closed")

	return nil
}
