package log

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// logEntry is the default implementation of the Entry interface
type logEntry struct {
	namespace   string
	level       int
	message     string
	isNoop      bool
	span        *trace.Span
	logFields   []zap.Field
	traceFields []trace.Attribute
	logger      *zap.Logger
}

func (l *logEntry) For(ctx context.Context) Entry {
	if span := trace.FromContext(ctx); span != nil {
		l.span = span
		l.WithField("TraceID", l.span.SpanContext().TraceID.String())
		l.WithField("SpanID", l.span.SpanContext().SpanID.String())
	}
	return l
}

func (l *logEntry) WithField(key string, value interface{}) Entry {
	logField, traceField := l.convertField(key, value)

	if l.logFields == nil {
		l.logFields = []zap.Field{logField}
	} else {
		l.logFields = append(l.logFields, logField)
	}

	if l.traceFields == nil {
		l.traceFields = []trace.Attribute{traceField}
	} else {
		l.traceFields = append(l.traceFields, traceField)
	}

	return l
}

func (l logEntry) Send() {
	if l.isNoop {
		return
	}

	if l.span != nil {
		l.span.Annotate(l.traceFields, l.message)
	}

	switch l.level {
	case LevelDebug:
		l.logger.Debug(l.message, l.logFields...)
	case LevelInfo:
		l.logger.Info(l.message, l.logFields...)
	case LevelWarning:
		l.logger.Warn(l.message, l.logFields...)
	case LevelError:
		l.logger.Error(l.message, l.logFields...)
	}
}

func (l logEntry) convertField(key string, value interface{}) (logField zap.Field, traceField trace.Attribute) {
	switch value.(type) {
	case []byte:
		v := string(value.([]byte))
		traceField = trace.StringAttribute(key, v)
		logField = zap.String(key, v)
	case []rune:
		v := string(value.([]rune))
		traceField = trace.StringAttribute(key, v)
		logField = zap.String(key, v)
	case string:
		v := value.(string)
		traceField = trace.StringAttribute(key, v)
		logField = zap.String(key, v)
	case bool:
		v := value.(bool)
		traceField = trace.BoolAttribute(key, v)
		logField = zap.Bool(key, v)
	case int:
		v := int64(value.(int))
		traceField = trace.Int64Attribute(key, v)
		logField = zap.Int64(key, v)
	case int32:
		v := int64(value.(int32))
		traceField = trace.Int64Attribute(key, v)
		logField = zap.Int64(key, v)
	case int64:
		v := value.(int64)
		traceField = trace.Int64Attribute(key, v)
		logField = zap.Int64(key, v)
	case uint:
		v := int64(value.(uint))
		traceField = trace.Int64Attribute(key, v)
		logField = zap.Int64(key, v)
	case uint32:
		v := int64(value.(uint32))
		traceField = trace.Int64Attribute(key, v)
		logField = zap.Int64(key, v)
	case uint64:
		v := int64(value.(uint64))
		traceField = trace.Int64Attribute(key, v)
		logField = zap.Int64(key, v)
	case float32:
		v := float64(value.(float32))
		traceField = trace.Float64Attribute(key, v)
		logField = zap.Float64(key, v)
	case float64:
		v := value.(float64)
		traceField = trace.Float64Attribute(key, v)
		logField = zap.Float64(key, v)
	default:
		v := fmt.Sprint(value)
		traceField = trace.StringAttribute(key, v)
		logField = zap.String(key, v)
	}

	return
}
