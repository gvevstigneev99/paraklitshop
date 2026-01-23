package service

import (
	"paraklitshop/internal/model"
	"paraklitshop/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAllProducts()
}
