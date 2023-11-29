package logger

import (
	"fmt"
	"os"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env, appName, appVersion, logPath string) (*zap.Logger, error) {
	cores, err := getCores(logPath)
	if err != nil {
		return nil, err
	}

	return zap.New(
		zapcore.NewTee(cores...),
		zap.Fields(zap.String(string(semconv.ServiceNameKey), appName)),
		zap.Fields(zap.String(string(semconv.ServiceVersionKey), appVersion)),
		zap.Fields(zap.String(string(semconv.DeploymentEnvironmentKey), env)),
	), nil
}

func getCores(logPath string) ([]zapcore.Core, error) {
	fileCore, err := getFileCore(logPath)
	if err != nil {
		return nil, err
	}

	consoleCore := getConsoleCore()

	return []zapcore.Core{
		consoleCore,
		fileCore,
	}, nil
}

func getFileCore(logPath string) (zapcore.Core, error) {
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("unable to open log file: %w", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writeSyncer := zapcore.Lock(file)
	lvlEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	return zapcore.NewCore(encoder, writeSyncer, lvlEnabler), nil
}

func getConsoleCore() zapcore.Core {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	writeSyncer := zapcore.Lock(os.Stdout)
	lvlEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	return zapcore.NewCore(encoder, writeSyncer, lvlEnabler)
}
