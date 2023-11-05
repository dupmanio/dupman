package processor

import (
	"encoding/json"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/scanner/broker"
	"github.com/dupmanio/dupman/packages/scanner/config"
	"github.com/dupmanio/dupman/packages/scanner/fetcher"
	"github.com/dupmanio/dupman/packages/scanner/model"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Processor struct {
	logger            *zap.Logger
	config            *config.Config
	broker            *broker.RabbitMQ
	fetcher           *fetcher.Fetcher
	dupmanCredentials credentials.Provider
	dupmanAPIService  *service.DupmanAPIService
}

func NewProcessor(
	logger *zap.Logger,
	config *config.Config,
	broker *broker.RabbitMQ,
	fetcher *fetcher.Fetcher,
	dupmanAPIService *service.DupmanAPIService,
) (*Processor, error) {
	cred, err := credentials.NewClientCredentials(
		config.DupmanConfig.ClientID,
		config.DupmanConfig.ClientSecret,
		[]string{},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman credentials provider: %w", err)
	}

	return &Processor{
		logger:            logger,
		config:            config,
		broker:            broker,
		fetcher:           fetcher,
		dupmanCredentials: cred,
		dupmanAPIService:  dupmanAPIService,
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
		zap.String("exchange", delivery.Exchange),
		zap.String("routingKey", delivery.RoutingKey),
	)

	err := proc.processScanning(delivery)
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

func (proc *Processor) processScanning(delivery amqp.Delivery) error {
	var message dto.ScanWebsiteMessage

	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		return fmt.Errorf("unable to unmarshal message: %w", err)
	}

	status := dto.Status{
		State: dto.StatusStateUpToDated,
	}

	updates, err := proc.fetcher.Fetch(message.WebsiteURL, message.WebsiteID, message.WebsiteToken)
	if err != nil {
		proc.logger.Error(
			"Unable to fetch Website Updates",
			zap.String("websiteID", message.WebsiteID.String()),
			zap.Error(err),
		)

		status.State = dto.StatusStateScanningFailed
		status.Info = err.Error()
	}

	if len(updates) != 0 {
		status.State = dto.StatusStateNeedsUpdate
	}

	if err = proc.updateWebsiteStatus(message.WebsiteID, status, updates); err != nil {
		return fmt.Errorf("unable to create Website Updates: %w", err)
	}

	return nil
}

func (proc *Processor) updateWebsiteStatus(websiteID uuid.UUID, status dto.Status, updatesRaw []model.Update) error {
	var updates dto.Updates

	_ = copier.Copy(&updates, &updatesRaw)

	proc.logger.Info(
		"Updating website status",
		zap.String("websiteID", websiteID.String()),
	)

	if err := proc.dupmanAPIService.CreateSession(proc.dupmanCredentials); err != nil {
		return fmt.Errorf("unable to create dupman session: %w", err)
	}

	_, err := proc.dupmanAPIService.SystemSvc.UpdateWebsiteStatus(websiteID, &status, &updates)
	if err != nil {
		return fmt.Errorf("unable to create Website Updates: %w", err)
	}

	proc.logger.Info(
		"Website status have been updated",
		zap.String("websiteID", websiteID.String()),
	)

	return nil
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
