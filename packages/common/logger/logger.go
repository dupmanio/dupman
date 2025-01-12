package logger

import (
	"os"

	"github.com/dupmanio/dupman/packages/common/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env string, ot *otel.OTel) (*zap.Logger, error) {
	return zap.New(
		zapcore.NewTee(
			zapcore.NewCore(getEncoder(env), zapcore.AddSync(os.Stdout), getLevelEnablerFunc(env)),
			ot.GetZapCore(),
		),
	), nil
}

func getEncoder(env string) zapcore.Encoder {
	if env == "prod" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

func getLevelEnablerFunc(env string) zapcore.Level {
	if env == "prod" {
		return zapcore.InfoLevel
	}

	return zapcore.DebugLevel
}
