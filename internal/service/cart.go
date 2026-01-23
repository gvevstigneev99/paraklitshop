package service

import (
	"errors"
	"paraklitshop/internal/repository"
)

type CartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) *CartService {
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
	return s.cartRepo.Add(userID, productID, qty)
}

func (s *CartService) GetCart(userID int) (map[int]int, error) {
	return s.cartRepo.Get(userID)
}

func (s *CartService) ClearCart(userID int) error {
	return s.cartRepo.Clear(userID)
}

func (s *CartService) RemoveItem(userID, productID int) error {
	return s.cartRepo.Remove(userID, productID)
}
