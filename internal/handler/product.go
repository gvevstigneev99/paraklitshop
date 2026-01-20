package handler

import (
	"paraklitshop/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) List(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch products"})
	}
	return c.JSON(products)
}
