package handler

import (
	"github.com/gofiber/fiber/v2"
)

// Health godoc
// @Summary Проверка состояния сервиса
// @Description Возвращает статус здоровья сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func Health() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "Service is healthy",
		})
	}
}
