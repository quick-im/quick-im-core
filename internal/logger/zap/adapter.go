package zap

import (
	"go.uber.org/zap"
)

type ZapLoggerAdapter struct {
	logger *zap.SugaredLogger
}

func NewZapLoggerAdapter(logger *zap.Logger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{logger: logger.Sugar()}
}

func (l *ZapLoggerAdapter) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *ZapLoggerAdapter) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *ZapLoggerAdapter) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *ZapLoggerAdapter) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *ZapLoggerAdapter) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *ZapLoggerAdapter) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
