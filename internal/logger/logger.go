package logger

import (
	"golang.org/x/exp/slog"

	"os"
)

func New(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case "local":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)

	// делаем его глобальным для стандартных вызовов slog.*
	slog.SetDefault(logger)

	return logger
}
