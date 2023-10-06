package broker

import (
	"fmt"
	"net"

	"github.com/dupmanio/dupman/packages/notifier/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	config     *config.Config
	logger     *zap.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewRabbitMQ(config *config.Config, logger *zap.Logger) (*RabbitMQ, error) {
	return &RabbitMQ{
		logger: logger,
		config: config,
	}, nil
}

func (brk *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	if err := brk.setup(); err != nil {
		return nil, fmt.Errorf("unable to setup RabbitMQ: %w", err)
	}

	messages, err := brk.channel.Consume(brk.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to start conssuming: %w", err)
	}

	return messages, nil
}

func (brk *RabbitMQ) setup() error {
	var err error

	url := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		brk.config.RabbitMQ.User,
		brk.config.RabbitMQ.Password,
		net.JoinHostPort(brk.config.RabbitMQ.Host, brk.config.RabbitMQ.Port),
	)

	brk.connection, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	brk.channel, err = brk.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	err = brk.channel.Qos(brk.config.Worker.PrefetchCount, brk.config.Worker.PrefetchSize, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	brk.queue, err = brk.channel.QueueDeclare(brk.config.RabbitMQ.QueueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
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
