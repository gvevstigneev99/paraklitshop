package service

import (
	"paraklitshop/internal/model"
	"paraklitshop/internal/repository/postgres"
)

type ProductService struct {
	repo *postgres.ProductRepository
}

func NewProductService(repo *postgres.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAllProducts()
}
