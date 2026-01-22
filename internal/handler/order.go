package handler

import (
	"paraklitshop/internal/service"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order based on the user's cart
// @Tags orders
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security     BearerAuth
// @Router /api/buyer/orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userIDLocal := c.Locals("userID")
	userID, ok := userIDLocal.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := h.service.CreateOrder(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to place order"})
	}
	return c.JSON(fiber.Map{"message": "Order successfully created"})
}
