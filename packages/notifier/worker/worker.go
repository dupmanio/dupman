package worker

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/notifier/broker"
	"github.com/dupmanio/dupman/packages/notifier/processor"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(lc fx.Lifecycle, logger *zap.Logger, broker *broker.RabbitMQ, processor *processor.Processor) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Worker")

			messages, err := broker.Consume()
			if err != nil {
				return fmt.Errorf("unable to consume messages: %w", err)
			}

			go func() {
				for msg := range messages {
					logger.Info(
						"Received a message",
						zap.Uint8("priority", msg.Priority),
						zap.String("messageID", msg.MessageId),
						zap.String("userID", msg.UserId),
						zap.String("appID", msg.AppId),
						zap.Uint64("deliveryTag", msg.DeliveryTag),
						zap.String("routingKey", msg.RoutingKey),
					)

					if err = processor.Process(msg); err != nil {
						logger.Error(
							"Unable to Process message",
							zap.String("messageID", msg.MessageId),
							zap.Error(err),
						)
					}

					if err = msg.Ack(false); err != nil {
						logger.Error(
							"Unable to Acknowledge message",
							zap.String("messageID", msg.MessageId),
							zap.Error(err),
						)
					}
				}
			}()

			logger.Info("Worker has been started. Waiting for messages")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down Worker")

			if err := broker.Shutdown(); err != nil {
				return fmt.Errorf("unable to shutdown broker: %w", err)
			}

			return nil
		},
	})

	return nil
}
