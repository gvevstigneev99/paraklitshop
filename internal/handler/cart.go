package handler

import (
	"paraklitshop/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID := 1 // In real scenario, get from JWT or session
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return err
	}
	qty, err := strconv.Atoi(c.Params("qty"))
	if err != nil {
		return err
	}
	err = h.service.AddItem(userID, productID, qty)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Item added to cart"})
}

func (h *CartHandler) ViewCart(c *fiber.Ctx) error {
	userID := 1 // In real scenario, get from JWT or session
	cart, err := h.service.GetCart(userID)
	if err != nil {
		return err
	}
	return c.JSON(cart)
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := 1 // In real scenario, get from JWT or session
	err := h.service.ClearCart(userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Cart cleared"})
}

func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	userID := 1 // In real scenario, get from JWT or session
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return err
	}
	err = h.service.RemoveItem(userID, productID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}
