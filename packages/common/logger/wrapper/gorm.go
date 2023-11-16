package wrapper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormWrapper struct {
	logger *zap.Logger
	config gormlogger.Config
}

func NewGormWrapper(logger *zap.Logger) *GormWrapper {
	logger = logger.With(
		zap.String("component", "gorm"),
	)

	return &GormWrapper{
		logger: logger,
		config: gormlogger.Config{
			LogLevel: gormlogger.Warn,
		},
	}
}

func (wrap *GormWrapper) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *wrap
	newLogger.config.LogLevel = level

	return &newLogger
}

func (wrap *GormWrapper) Info(_ context.Context, msg string, data ...interface{}) {
	if wrap.config.LogLevel >= gormlogger.Info {
		wrap.logger.Info(fmt.Sprintf(msg, data...))
	}
}

func (wrap *GormWrapper) Warn(_ context.Context, msg string, data ...interface{}) {
	if wrap.config.LogLevel >= gormlogger.Warn {
		wrap.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (wrap *GormWrapper) Error(_ context.Context, msg string, data ...interface{}) {
	if wrap.config.LogLevel >= gormlogger.Error {
		wrap.logger.Error(fmt.Sprintf(msg, data...))
	}
}

func (wrap *GormWrapper) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if wrap.config.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	logFields := []zap.Field{
		zap.String("caller", utils.FileWithLineNum()),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.String("elapsed", elapsed.String()),
	}

	switch {
	case err != nil &&
		wrap.config.LogLevel >= gormlogger.Error &&
		(!errors.Is(err, gormlogger.ErrRecordNotFound) || !wrap.config.IgnoreRecordNotFoundError):
		wrap.logger.Error(
			"gorm error occurred",
			append(logFields, zap.Error(err))...,
		)
	case elapsed > wrap.config.SlowThreshold && wrap.config.SlowThreshold != 0 && wrap.config.LogLevel >= gormlogger.Warn:
		wrap.logger.Warn(
			"slow SQL",
			append(logFields, zap.Duration("duration", wrap.config.SlowThreshold))...,
		)
	case wrap.config.LogLevel == gormlogger.Info:
		wrap.logger.Debug(
			"sql tracing",
			logFields...,
		)
	}
}