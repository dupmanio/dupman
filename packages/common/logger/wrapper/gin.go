package wrapper

import (
	"regexp"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GinWrapper struct {
	logger *zap.Logger
}

func NewGinWrapper(logger *zap.Logger) *GinWrapper {
	logger = logger.With(
		zap.String("component", "gin"),
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
			zap.String("method", matches[1]),
			zap.String("route", matches[2]),
			zap.String("handler", matches[3]),
			zap.String("handlersNum", matches[4]),
		)

		return
	}

	wrap.logger.Debug(message)
}

func (wrap *GinWrapper) GetGinzapMiddleware() gin.HandlerFunc {
	// @todo: fix TraceID.
	return ginzap.GinzapWithConfig(
		wrap.logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			TraceID:    true,
		},
	)
}

func (wrap *GinWrapper) GetGinzapRecoveryMiddleware() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(wrap.logger, true)
}
