package log

import (
	"errors"

	"go.uber.org/zap"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
)

var (
	ErrNamespaceMissing = errors.New("logger namespace missing")
	ErrZapLoggerMissing = errors.New("logger zap config missing")
)

// StandardLogger is the default implementation of the Logger interface
type StandardLogger struct {
	namespace string
	isNoop    bool
	logger    *zap.Logger
}

func New(namespace string, zapLogger *zap.Logger, options ...Option) (*StandardLogger, error) {
	if namespace == "" {
		return nil, ErrNamespaceMissing
	}

	if zapLogger == nil {
		return nil, ErrZapLoggerMissing
	}

	l := &StandardLogger{
		namespace: namespace,
		isNoop:    false,
		logger:    zapLogger,
	}

	for _, opt := range options {
		opt.apply(l)
	}

	return l, nil
}

func NewDefault(namespace string, options ...Option) (*StandardLogger, error) {
	cfg := DefaultConfig(namespace)
	return New(namespace, cfg, options...)
}

func (s *StandardLogger) Debug(msg string) Entry {
	return s.newLogEntry(LevelDebug, msg)
}

func (s *StandardLogger) Info(msg string) Entry {
	return s.newLogEntry(LevelInfo, msg)
}

func (s *StandardLogger) Warning(msg string) Entry {
	return s.newLogEntry(LevelWarning, msg)
}

func (s *StandardLogger) Error(msg string) Entry {
	return s.newLogEntry(LevelError, msg)
}

func (s *StandardLogger) Close() error {
	return s.logger.Sync()
}

func (s *StandardLogger) newLogEntry(level int, msg string) Entry {
	return &logEntry{
		namespace: s.namespace,
		level:     level,
		message:   msg,
		isNoop:    s.isNoop,
		span:      nil,
		logger:    s.logger,
	}
}
