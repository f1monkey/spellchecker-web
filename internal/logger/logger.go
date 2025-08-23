package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxType struct{}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxType{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxType{}).(*slog.Logger); ok {
		return l
	}

	return noop
}

func New(
	appVersion string,
	level string,
) *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: levelFromString(level),
			},
		),
	)

	logger = logger.With("version", appVersion)

	return logger
}

func levelFromString(level string) slog.Leveler {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
