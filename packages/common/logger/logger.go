package logger

import (
	"os"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env, appName, appVersion string) (*zap.Logger, error) {
	core := zapcore.NewCore(
		getEncoder(env),
		zapcore.Lock(os.Stdout),
		getLevelEnablerFunc(env),
	)

	return zap.New(
		core,
		zap.Fields(zap.String(string(semconv.ServiceNameKey), appName)),
		zap.Fields(zap.String(string(semconv.ServiceVersionKey), appVersion)),
		zap.Fields(zap.String(string(semconv.DeploymentEnvironmentKey), env)),
	), nil
}

func getEncoder(env string) zapcore.Encoder {
	if env == "prod" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

func getLevelEnablerFunc(env string) zap.LevelEnablerFunc {
	minLevel := zapcore.DebugLevel
	if env == "prod" {
		minLevel = zapcore.InfoLevel
	}

	return func(lvl zapcore.Level) bool {
		return lvl >= minLevel
	}
}
