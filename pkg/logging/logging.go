package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type ctxKey struct{}

var logger *slog.Logger

// initialize the global logger with the given service name
func Init(service string) {
	level := parseLogLevel(os.Getenv("LOG_LEVEL"))

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger = slog.New(handler).With("service", service)
}

// returns a context that stores structured fields (e.g., request_id)
func WithContext(ctx context.Context, keyvals ...any) context.Context {
	fields := make(map[string]any)
	for i := 0; i < len(keyvals)-1; i += 2 {
		k, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		fields[k] = keyvals[i+1]
	}
	return context.WithValue(ctx, ctxKey{}, fields)
}

// extracts a logger with context fields
func From(ctx context.Context) *slog.Logger {
	log := logger
	if fields, ok := ctx.Value(ctxKey{}).(map[string]any); ok {
		for k, v := range fields {
			log = log.With(k, v)
		}
	}
	return log
}

// Shorthand logging methods
func Debug(msg string, keyvals ...any) {
	logger.Debug(msg, keyvals...)
}

func Warn(msg string, keyvals ...any) {
	logger.Warn(msg, keyvals...)
}

func Error(msg string, keyvals ...any) {
	logger.Error(msg, keyvals...)
}


// Helper to parse log level from environment
func parseLogLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
