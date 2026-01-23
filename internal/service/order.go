package service

import (
	"errors"
	"paraklitshop/internal/model"
	"paraklitshop/internal/repository"
	"time"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(userID int) error {
	cartItems, err := s.cartRepo.Get(userID)
	if err != nil {
		return err
	}

	if len(cartItems) == 0 {
		return errors.New("cart is empty")
	}

	var totalAmount float64

	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return err
	}
	for productID, qty := range cartItems {
		var productFound bool
		for _, p := range products {
			if p.ID == productID {
				totalAmount += p.Price * float64(qty)
				productFound = true
				break
			}
		}
		if !productFound {
			return errors.New("product not found")
		}
	}

	now := time.Now()
	order := model.Order{
		ID:         0,
		UserID:     userID,
		TotalPrice: totalAmount,
		Status:     "paid",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = s.orderRepo.Create(&order)
	if err != nil {
		return err
	}

	err = s.cartRepo.Clear(userID)
	if err != nil {
		return err
	}

	return nil
}
