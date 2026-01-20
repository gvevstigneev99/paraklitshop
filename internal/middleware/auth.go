package middleware

import (
	"github.com/gofiber/fiber/v2"
)

const (
	validToken = "secret123"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token != "Bearer "+validToken {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return c.Next()
	}
}
