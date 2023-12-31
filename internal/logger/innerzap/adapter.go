package innerzap

import (
	"go.uber.org/zap"
)

type ZapLoggerAdapter struct {
	logger *zap.Logger
}

func NewZapLoggerAdapter(logger *zap.Logger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{logger: logger.WithOptions(zap.AddCallerSkip(1))}
}

func (z *ZapLoggerAdapter) Debug(msg string, args ...string) {
	fields := make([]zap.Field, len(args))
	for i := range args {
		fields[i] = zap.String("", args[i])
	}
	z.logger.Debug(msg, fields...)
}

func (z *ZapLoggerAdapter) Info(msg string, args ...string) {
	fields := make([]zap.Field, len(args))
	for i := range args {
		fields[i] = zap.String("", args[i])
	}
	z.logger.Info(msg, fields...)
}

func (z *ZapLoggerAdapter) Warn(msg string, args ...string) {
	fields := make([]zap.Field, len(args))
	for i := range args {
		fields[i] = zap.String("", args[i])
	}
	z.logger.Warn(msg, fields...)
}

func (z *ZapLoggerAdapter) Error(msg string, args ...string) {
	fields := make([]zap.Field, len(args))
	for i := range args {
		fields[i] = zap.String("", args[i])
	}
	z.logger.Error(msg, fields...)
}

func (z *ZapLoggerAdapter) Panic(msg string, args ...string) {
	fields := make([]zap.Field, len(args))
	for i := range args {
		fields[i] = zap.String("", args[i])
	}
	z.logger.Panic(msg, fields...)
}

func (z *ZapLoggerAdapter) Fatal(msg string, args ...string) {
	fields := make([]zap.Field, len(args))
	for i := range args {
		fields[i] = zap.String("", args[i])
	}
	z.logger.Fatal(msg, fields...)
}
