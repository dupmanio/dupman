package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/otel"
	commonService "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/scanner/config"
	"github.com/dupmanio/dupman/packages/scanner/fetcher"
	"github.com/dupmanio/dupman/packages/scanner/model"
	"github.com/dupmanio/dupman/packages/scanner/service"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	amqp "github.com/rabbitmq/amqp091-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Processor struct {
	logger            *zap.Logger
	config            *config.Config
	messengerSvc      *service.MessengerService
	fetcher           *fetcher.Fetcher
	dupmanCredentials credentials.Provider
	dupmanAPIService  *commonService.DupmanAPIService
	ot                *otel.OTel
}

func NewProcessor(
	logger *zap.Logger,
	config *config.Config,
	messengerSvc *service.MessengerService,
	fetcher *fetcher.Fetcher,
	dupmanAPIService *commonService.DupmanAPIService,
	ot *otel.OTel,
) (*Processor, error) {
	cred, err := credentials.NewClientCredentials(
		config.Dupman.ClientID,
		config.Dupman.ClientSecret,
		[]string{},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman credentials provider: %w", err)
	}

	return &Processor{
		logger:            logger,
		config:            config,
		messengerSvc:      messengerSvc,
		fetcher:           fetcher,
		dupmanCredentials: cred,
		dupmanAPIService:  dupmanAPIService,
		ot:                ot,
	}, nil
}

func (proc *Processor) Process() error {
	messages, err := proc.messengerSvc.Consume()
	if err != nil {
		return fmt.Errorf("unable to consume messages: %w", err)
	}

	for msg := range messages {
		proc.processMessage(msg)
	}

	return nil
}

func (proc *Processor) processMessage(delivery amqp.Delivery) {
	ctx := proc.getOTelContext(delivery)

	ctx, span := proc.ot.Tracer.Start(
		ctx,
		fmt.Sprintf("%s receive", proc.config.Worker.QueueName),
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingOperationReceive,
			semconv.MessagingMessageID(delivery.MessageId),

			// semconv.MessagingSystem("rabbitmq"),
			// semconv.MessagingClientID(proc.config.RabbitMQ),
			// semconv.ServerAddress(proc.config.RabbitMQ.Host),
		),
	)
	defer span.End()

	proc.ot.LogInfoEvent(
		ctx,
		"Received a message",
		semconv.MessagingMessageID(delivery.MessageId),
	)

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

	successfullyProcessed := true
	if err := proc.processScanning(delivery); err != nil {
		successfullyProcessed = false

		proc.logger.Error(
			"Unable to Process message",
			zap.String("messageID", delivery.MessageId),
			zap.Error(err),
		)
	}

	proc.messengerSvc.AcknowledgeMessage(successfullyProcessed, delivery)
}

func (proc *Processor) getOTelContext(delivery amqp.Delivery) context.Context {
	ctx := context.Background()
	traceID, _ := trace.TraceIDFromHex(proc.getDeliveryStringHeader(delivery, "trace_id"))
	spanID, _ := trace.SpanIDFromHex(proc.getDeliveryStringHeader(delivery, "span_id"))

	if traceID.IsValid() && spanID.IsValid() {
		ctx = trace.ContextWithRemoteSpanContext(
			ctx,
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				// SpanID:  spanID,
			}),
		)
	}

	return ctx
}

func (proc *Processor) getDeliveryStringHeader(delivery amqp.Delivery, headerName string) string {
	if headerRaw := delivery.Headers[headerName]; headerRaw != nil {
		if headerString, ok := headerRaw.(string); ok {
			return headerString
		}
	}

	return ""
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
