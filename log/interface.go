package log

import "context"

type Logger interface {
	Debug(string) Entry   // Init new LogEntry with level Debug
	Info(string) Entry    // Init new LogEntry with level Info
	Warning(string) Entry // Init new LogEntry with level Warning
	Error(string) Entry   // Init new LogEntry with level Error
	Panic(string) Entry   // Init new LogEntry with level Panic (will call panic())
}

type Entry interface {
	// Attach LogEntry to Context (for TraceID propagation)
	For(context.Context) Entry

	// Add custom data to LogEntry
	WithField(key string, value interface{}) Entry

	// Send LogEntry to stdout
	Send()
}
