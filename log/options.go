package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Options to pass into the Logger
type Option interface {
	apply(*StandardLogger)
}

type optionFn func(*StandardLogger)

func (o optionFn) apply(l *StandardLogger) {
	o(l)
}

func Noop() Option {
	return optionFn(func(input *StandardLogger) {
		input.isNoop = true
	})
}

func DefaultConfig(namespace string) *zap.Logger {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	conf.EncoderConfig.FunctionKey = "func"
	conf.OutputPaths = []string{"stdout"}
	conf.ErrorOutputPaths = []string{"stderr"}
	conf.Development = false
	zapLogger, _ := conf.Build(
		zap.Fields(zap.String("namespace", namespace)),
		zap.AddStacktrace(zapcore.PanicLevel),
		zap.AddCallerSkip(1),
	)

	return zapLogger
}

func ObservedConfig(namespace string) (*zap.Logger, *observer.ObservedLogs) {
	wrappedCore, obs := observer.New(zap.NewAtomicLevelAt(zap.InfoLevel))

	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	conf.EncoderConfig.FunctionKey = "func"
	conf.OutputPaths = []string{"stdout"}
	conf.ErrorOutputPaths = []string{"stderr"}
	conf.Development = false
	zapLogger, _ := conf.Build(
		zap.Fields(zap.String("namespace", namespace)),
		zap.AddStacktrace(zapcore.PanicLevel),
		zap.AddCallerSkip(1),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return wrappedCore
		}),
	)

	return zapLogger, obs
}
