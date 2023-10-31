package processor

import (
	"encoding/json"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/notifier/broker"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"github.com/dupmanio/dupman/packages/notifier/deliverer"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/email"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/notify"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Processor struct {
	logger            *zap.Logger
	config            *config.Config
	broker            *broker.RabbitMQ
	dupmanCredentials credentials.Provider
	dupmanAPIService  *service.DupmanAPIService
	deliverers        []deliverer.Deliverer
}

func NewProcessor(
	logger *zap.Logger,
	config *config.Config,
	broker *broker.RabbitMQ,
	dupmanAPIService *service.DupmanAPIService,
	emailDeliverer *email.Deliverer,
	notifyDeliverer *notify.Deliverer,
) (*Processor, error) {
	cred, err := credentials.NewClientCredentials(
		config.DupmanAPIService.ClientID,
		config.DupmanAPIService.ClientSecret,
		[]string{"user:get_contact_info"},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman credentials provider: %w", err)
	}

	return &Processor{
		logger:            logger,
		config:            config,
		broker:            broker,
		dupmanCredentials: cred,
		dupmanAPIService:  dupmanAPIService,
		deliverers: []deliverer.Deliverer{
			emailDeliverer,
			notifyDeliverer,
		},
	}, nil
}

func (proc *Processor) Process() error {
	messages, err := proc.broker.Consume()
	if err != nil {
		return fmt.Errorf("unable to consume messages: %w", err)
	}

	for msg := range messages {
		proc.processMessage(msg)
	}

	return nil
}

func (proc *Processor) processMessage(delivery amqp.Delivery) {
	proc.logger.Info(
		"Received a message",
		zap.String("messageID", delivery.MessageId),
		zap.Uint8("priority", delivery.Priority),
		zap.String("userID", delivery.UserId),
		zap.String("appID", delivery.AppId),
		zap.Uint64("deliveryTag", delivery.DeliveryTag),
		zap.String("routingKey", delivery.RoutingKey),
	)

	err := proc.processDelivery(delivery)
	if err != nil {
		proc.logger.Error(
			"Unable to Process message",
			zap.String("messageID", delivery.MessageId),
			zap.Error(err),
		)
	}

	if err == nil || (err != nil && proc.isLastDeliveryAttempt(delivery)) {
		if err = delivery.Ack(false); err != nil {
			proc.logger.Error(
				"Unable to Ack message",
				zap.String("messageID", delivery.MessageId),
				zap.Error(err),
			)
		}
	} else {
		if err = delivery.Nack(false, false); err != nil {
			proc.logger.Error(
				"Unable to Nack message",
				zap.String("messageID", delivery.MessageId),
				zap.Error(err),
			)
		}
	}
}

func (proc *Processor) processDelivery(delivery amqp.Delivery) error {
	var message dto.NotificationMessage

	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		return fmt.Errorf("unable to unmarshal message: %w", err)
	}

	contactInfo, err := proc.getUserContactInfo(message.UserID)
	if err != nil {
		return fmt.Errorf("unable to get user contact info: %w", err)
	}

	for _, delivererInstance := range proc.deliverers {
		go proc.attemptNotificationDelivery(delivererInstance, contactInfo, delivery.MessageId, message)
	}

	return nil
}

func (proc *Processor) getUserContactInfo(userID uuid.UUID) (*dto.ContactInfo, error) {
	if err := proc.dupmanAPIService.CreateSession(proc.dupmanCredentials); err != nil {
		return nil, fmt.Errorf("unable to create dupman session: %w", err)
	}

	proc.logger.Info(
		"Fetching user contact info",
		zap.String("userID", userID.String()),
	)

	info, err := proc.dupmanAPIService.UserSvc.GetContactInfo(userID)
	if err != nil {
		proc.logger.Error(
			"Unable to fetch user contact info",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)

		return nil, fmt.Errorf("unable to fetch user contact info: %w", err)
	}

	return info, nil
}

func (proc *Processor) attemptNotificationDelivery(
	delivererInstance deliverer.Deliverer,
	contactInfo *dto.ContactInfo,
	messageID string,
	message dto.NotificationMessage,
) {
	logFields := []zap.Field{
		zap.String("messageID", messageID),
		zap.String("userID", message.UserID.String()),
		zap.String("messageType", message.Type),
		zap.String("deliverer", delivererInstance.Name()),
	}

	proc.logger.Info("Starting delivering notification", logFields...)

	if err := delivererInstance.Deliver(message, contactInfo); err == nil {
		proc.logger.Info("Notification has been delivered", logFields...)
	} else {
		proc.logger.Error("Unable to deliver notification", append(logFields, zap.Error(err))...)
	}
}

func (proc *Processor) isLastDeliveryAttempt(delivery amqp.Delivery) bool {
	xDeath, ok := delivery.Headers["x-death"].([]interface{})
	if !ok {
		return false
	}

	count, ok := xDeath[0].(amqp.Table)["count"].(int64)
	if !ok {
		return false
	}

	return int(count) >= proc.config.Worker.RetryAttempts
}

func (proc *Processor) Shutdown() error {
	if err := proc.broker.Shutdown(); err != nil {
		return fmt.Errorf("unable to shutdown broker: %w", err)
	}

	return nil
}
