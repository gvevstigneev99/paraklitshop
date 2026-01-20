package logger

import (
	"log"

	"golang.org/x/exp/slog"
)

func New(env string) *slog.Logger {
	var level slog.Level
	if env == "local" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}
