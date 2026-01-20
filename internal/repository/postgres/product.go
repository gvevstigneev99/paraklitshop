package postgres

import "paraklitshop/internal/model"

type ProductRepository struct {
	products []model.Product
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: []model.Product{
			{ID: 1, Title: "Hoddie Paraklit", Description: "Description 1", Price: 4000.0, SellerID: 1},
			{ID: 2, Title: "T-short Paraklit", Description: "Description 2", Price: 3000.0, SellerID: 2},
			{ID: 3, Title: "Disk 1", Description: "Description 3", Price: 300.0, SellerID: 1},
		},
	}
}

func (r *ProductRepository) GetAllProducts() ([]model.Product, error) {
	return r.products, nil
}
