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

func userIDFromLocals(c *fiber.Ctx) (int, bool) {
	v := c.Locals("userID")
	userID, ok := v.(int)
	return userID, ok
}

// AddToCart godoc
// @Summary Добавить товар в корзину
// @Description Добавляет указанный товар в корзину пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param productId path int true "ID товара"
// @Param qty path int true "Количество"
// @Security     BearerAuth
// @Router /cart/add/{productId} [post]
func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID, ok := userIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
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

// ViewCart godoc
// @Summary Просмотр корзины
// @Description Возвращает содержимое корзины пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /cart [get]
func (h *CartHandler) ViewCart(c *fiber.Ctx) error {
	userID, ok := userIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	cart, err := h.service.GetCart(userID)
	if err != nil {
		return err
	}
	return c.JSON(cart)
}

// ClearCart godoc
// @Summary Очистить корзину
// @Description Удаляет все товары из корзины пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /cart/clear [delete]
func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID, ok := userIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	err := h.service.ClearCart(userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Cart cleared"})
}

// RemoveFromCart godoc
// @Summary Удалить товар из корзины
// @Description Удаляет указанный товар из корзины пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param productId path int true "ID товара"
// @Security     BearerAuth
// @Router /cart/remove/{productId} [delete]
func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	userID, ok := userIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
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
