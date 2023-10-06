package processor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"github.com/dupmanio/dupman/packages/notifier/deliverer"
	"github.com/dupmanio/dupman/packages/notifier/deliverer/email"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Processor struct {
	logger     *zap.Logger
	config     *config.Config
	deliverers []deliverer.Deliverer
}

func NewProcessor(logger *zap.Logger, config *config.Config, emailDeliverer *email.Deliverer) (*Processor, error) {
	return &Processor{
		logger: logger,
		config: config,
		deliverers: []deliverer.Deliverer{
			emailDeliverer,
		},
	}, nil
}

func (proc *Processor) Process(delivery amqp.Delivery) error {
	var message dto.NotificationMessage

	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		return fmt.Errorf("unable to unmarshal message: %w", err)
	}

	contactInfo, err := proc.getUserContactInfo(message.UserID)
	if err != nil {
		return fmt.Errorf("unable to get user contact info: %w", err)
	}

	proc.deliverNotification(delivery, message, contactInfo)

	return nil
}

func (proc *Processor) getUserContactInfo(userID uuid.UUID) (dto.UserContactInfo, error) {
	// @todo: replace with actual implementation.

	return dto.UserContactInfo{
		Email: "user1@ema.il",
	}, nil
}

func (proc *Processor) deliverNotification(
	delivery amqp.Delivery,
	message dto.NotificationMessage,
	contactInfo dto.UserContactInfo,
) {
	for _, delivererInstance := range proc.deliverers {
		go func(delivererInstance deliverer.Deliverer) {
			proc.logger.Info(
				"Starting delivering notification",
				zap.String("messageID", delivery.MessageId),
				zap.String("userID", message.UserID.String()),
				zap.String("messageType", message.Type),
				zap.String("deliverer", delivererInstance.Name()),
			)

			if proc.tryToDeliverNotificationUsingSingleDeliverer(delivererInstance, contactInfo, message, delivery) {
				proc.logger.Info(
					"Notification has been delivered",
					zap.String("messageID", delivery.MessageId),
					zap.String("userID", message.UserID.String()),
					zap.String("messageType", message.Type),
					zap.String("deliverer", delivererInstance.Name()),
				)
			}
		}(delivererInstance)
	}
}

func (proc *Processor) tryToDeliverNotificationUsingSingleDeliverer(
	delivererInstance deliverer.Deliverer,
	contactInfo dto.UserContactInfo,
	message dto.NotificationMessage,
	delivery amqp.Delivery,
) bool {
	for retryAttempt := 1; retryAttempt <= proc.config.Deliverer.Retries; retryAttempt++ {
		err := delivererInstance.Deliver(message, contactInfo)
		if err == nil {
			return true
		}

		proc.logger.Error(
			"Notification delivery attempt has failed",
			zap.String("messageID", delivery.MessageId),
			zap.String("userID", message.UserID.String()),
			zap.String("messageType", message.Type),
			zap.String("deliverer", delivererInstance.Name()),
			zap.Int("retryAttempt", retryAttempt),
			zap.Error(err),
		)

		// If this is a last attempt, do not sleep.
		if retryAttempt < proc.config.Deliverer.Retries {
			time.Sleep(time.Duration(retryAttempt) * time.Second)
		}
	}

	return false
}
