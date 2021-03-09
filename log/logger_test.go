package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
	"go.uber.org/zap/zapcore"
)

func TestNewDefault(t *testing.T) {
	test := assert.New(t)

	l, err := NewDefault("test-1")
	ctx, span := trace.StartSpan(context.Background(), "test-context-1")
	defer span.End()
	if test.Nil(err) {
		l.Info("hello there").
			For(ctx).
			WithField("one", 1).
			WithField("two", "dos").
			WithField("three", true).
			WithField("four", 4.21).
			WithField("five", []byte("what")).
			WithField("six", []string{"A", "B", "C"}).
			WithField("seven", []rune{'d', 'e', 'f'}).
			Send()
	}
}

func TestNewObserved(t *testing.T) {
	test := assert.New(t)

	cfg, obs := ObservedConfig("test-1")
	l, err := New("test-1", cfg)

	ctx, span := trace.StartSpan(context.Background(), "test-context-1")
	defer span.End()

	expected := map[string]interface{}{
		"one":     int64(1),
		"two":     "dos",
		"three":   true,
		"four":    4.21,
		"five":    "what",
		"six":     "[A B C]",
		"seven":   "def",
		"TraceID": span.SpanContext().TraceID.String(),
		"SpanID":  span.SpanContext().SpanID.String(),
	}

	if test.Nil(err) {
		l.Info("hello there").
			For(ctx).
			WithField("one", 1).
			WithField("two", "dos").
			WithField("three", true).
			WithField("four", 4.21).
			WithField("five", []byte("what")).
			WithField("six", []string{"A", "B", "C"}).
			WithField("seven", []rune{'d', 'e', 'f'}).
			Send()

		if test.Equal(1, obs.Len()) {
			encoder := zapcore.NewMapObjectEncoder()
			for _, field := range obs.All()[0].Context {
				field.AddTo(encoder)
			}

			test.Equal(expected, encoder.Fields)
		}
	}
}
