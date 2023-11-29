package wrapper

import (
	"regexp"
	"strings"
	"time"

	"github.com/dupmanio/dupman/packages/common/otel"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GinWrapper struct {
	logger *zap.Logger
}

func NewGinWrapper(logger *zap.Logger) *GinWrapper {
	logger = logger.With(
		zap.String(string(otel.ComponentKey), "gin"),
	)

	return &GinWrapper{
		logger: logger,
	}
}

func (wrap *GinWrapper) Write(data []byte) (int, error) {
	message := strings.TrimSuffix(string(data), "\n")

	switch {
	case strings.HasPrefix(message, "[GIN-debug]"):
		wrap.logGinDebug(strings.TrimPrefix(message, "[GIN-debug] "))
	case strings.HasPrefix(message, "[GIN-error]"):
		wrap.logger.Error(strings.TrimPrefix(message, "[GIN-debug] "))
	default:
		wrap.logger.Info(message)
	}

	return len(data), nil
}

func (wrap *GinWrapper) logGinDebug(message string) {
	re := regexp.MustCompile(`^(\w+)\s+(\S+)\s+-->\s+(\S+)\s+\((\d+)\s+handlers\)$`)
	if matches := re.FindStringSubmatch(message); len(matches) > 0 {
		wrap.logger.Debug(
			"New handler added",
			zap.String(string(semconv.HTTPMethodKey), matches[1]),
			zap.String(string(semconv.HTTPRouteKey), matches[2]),
			zap.String(string(otel.GinHandlerKey), matches[3]),
			zap.String(string(otel.GinHandlerNumKey), matches[4]),
		)

		return
	}

	wrap.logger.Debug(message)
}

func (wrap *GinWrapper) GetGinzapMiddleware() gin.HandlerFunc {
	return ginzap.GinzapWithConfig(
		wrap.logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			Context: func(c *gin.Context) []zapcore.Field {
				var fields []zapcore.Field

				if spanContext := trace.SpanFromContext(c.Request.Context()).SpanContext(); spanContext.IsValid() {
					if spanContext.TraceID().IsValid() {
						fields = append(fields, zap.String(string(otel.TraceIDKey), spanContext.TraceID().String()))
					}

					if spanContext.SpanID().IsValid() {
						fields = append(fields, zap.String(string(otel.SpanIDKey), spanContext.SpanID().String()))
					}
				}

				return fields
			},
		},
	)
}

func (wrap *GinWrapper) GetGinzapRecoveryMiddleware() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(wrap.logger, true)
}
