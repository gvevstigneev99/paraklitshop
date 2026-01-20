package handler

import (
	"github.com/gofiber/fiber/v2"
	"paraklitshop/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := 1 // In real scenario, get from JWT or session
	if err := h.service.CreateOrder(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to place order"})
	}
	return c.JSON(fiber.Map{"message": "Order successfully created"})
}
