package service

import (
	"paraklitshop/internal/repository/postgres"
	"paraklitshop/internal/repository/redis"

	"errors"
)

type CartService struct {
	cartRepo    *redis.CartRepository
	productRepo *postgres.ProductRepository
}

func NewCartService(cartRepo *redis.CartRepository, productRepo *postgres.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) AddItem(userID, productID, qty int) error {
	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return err
	}
	var productExists bool
	for _, p := range products {
		if p.ID == productID {
			productExists = true
			break
		}
	}
	if !productExists {
		return errors.New("product does not exist")
	}
	return s.cartRepo.AddItem(userID, productID, qty)
}

func (s *CartService) GetCart(userID int) (map[int]int, error) {
	return s.cartRepo.GetCart(userID)
}

func (s *CartService) ClearCart(userID int) error {
	return s.cartRepo.ClearCart(userID)
}

func (s *CartService) RemoveItem(userID, productID int) error {
	return s.cartRepo.RemoveItem(userID, productID)
}
