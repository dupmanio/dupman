package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type RabbitMQConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	AppID    string
}

type RabbitMQ struct {
	ot         *otel.OTel
	config     *RabbitMQConfig
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(ot *otel.OTel, config *RabbitMQConfig) (*RabbitMQ, error) {
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
		ot:         ot,
		config:     config,
		Connection: connection,
		Channel:    channel,
	}, nil
}

func (brk *RabbitMQ) PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message any) error {
	const timeout = 5 * time.Second

	messageID := uuid.New().String()

	ctx, span := brk.ot.Tracer.Start(
		ctx,
		fmt.Sprintf("%s.%s publish", exchangeName, routingKey),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(brk.getCommonOTelAttributes()...),
		trace.WithAttributes(
			semconv.MessagingOperationPublish,
			semconv.MessagingMessageID(messageID),
			semconv.MessagingDestinationName(exchangeName),
			semconv.MessagingRabbitmqDestinationRoutingKey(routingKey),
		),
	)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// @todo: add some logging.
	err = brk.Channel.PublishWithContext(ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{
				"trace_id": span.SpanContext().TraceID().String(),
				"span_id":  span.SpanContext().SpanID().String(),
			},
			MessageId:    messageID,
			UserId:       brk.config.User,
			AppId:        brk.config.AppID,
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonMessage,
		})
	if err != nil {
		return fmt.Errorf("failed to publish to channel: %w", err)
	}

	messagesPublished, _ := brk.ot.Meter.Int64Counter(
		"go.rabbitmq.message.published",
		metric.WithDescription("The total number of messages published"),
	)
	messagesPublished.Add(ctx, 1)

	return nil
}

func (brk *RabbitMQ) getCommonOTelAttributes() []attribute.KeyValue {
	fields := []attribute.KeyValue{
		semconv.MessagingSystem("rabbitmq"),
		semconv.MessagingClientID(brk.config.AppID),
		semconv.ServerAddress(brk.config.Host),
	}

	if port, err := strconv.Atoi(brk.config.Port); err == nil {
		fields = append(fields, semconv.ServerPort(port))
	}

	return fields
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
