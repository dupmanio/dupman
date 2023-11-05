package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	AppID    string
}

type RabbitMQ struct {
	config     *RabbitMQConfig
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(config *RabbitMQConfig) (*RabbitMQ, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		config.User,
		config.Password,
		net.JoinHostPort(config.Host, config.Port),
	)

	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a Channel: %w", err)
	}

	return &RabbitMQ{
		config:     config,
		Connection: connection,
		Channel:    channel,
	}, nil
}

func (brk *RabbitMQ) PublishToExchange(exchangeName string, routingKey string, message any) error {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// @todo: Implement ack checking and retry functionality.
	// pubAck, pubNack := brk.Channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))
	// select {
	// case <-pubAck:
	//	 return nil
	// case <-pubNack:
	//	 return domainErrors.ErrUnableToPublishMessage
	// case <-time.After(5 * time.Second):
	//	 return domainErrors.ErrUnableToPublishMessage
	// }
	// @todo: add some logging.
	err = brk.Channel.PublishWithContext(ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			MessageId:    uuid.New().String(),
			UserId:       brk.config.User,
			AppId:        brk.config.AppID,
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonMessage,
		})
	if err != nil {
		return fmt.Errorf("failed to publish to channel: %w", err)
	}

	return nil
}

func (brk *RabbitMQ) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {
	messages, err := brk.Channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to consume queue: %w", err)
	}

	return messages, nil
}

func (brk *RabbitMQ) GetDeathCount(delivery amqp.Delivery) int64 {
	xDeath, ok := delivery.Headers["x-death"].([]interface{})
	if !ok {
		return 0
	}

	count, ok := xDeath[0].(amqp.Table)["count"].(int64)
	if !ok {
		return 0
	}

	return count
}

func (brk *RabbitMQ) Close() error {
	if err := brk.Channel.Close(); err != nil {
		return fmt.Errorf("unable to close channel: %w", err)
	}

	if err := brk.Connection.Close(); err != nil {
		return fmt.Errorf("unable to close connection: %w", err)
	}

	return nil
}
