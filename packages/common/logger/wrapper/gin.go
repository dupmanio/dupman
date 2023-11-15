package wrapper

import (
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
		wrap.logger.Debug(strings.TrimPrefix(message, "[GIN-debug]"))
	case strings.HasPrefix(message, "[GIN-error]"):
		wrap.logger.Error(strings.TrimPrefix(message, "[GIN-debug]"))
	default:
		wrap.logger.Info(message)
	}

	return len(data), nil
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
