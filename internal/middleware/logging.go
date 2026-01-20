package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func LoggingMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		logger.Info("request completed", slog.String("method", c.Method()), slog.String("path", c.Path()), slog.Duration("duration", duration))
		return err
	}
}
