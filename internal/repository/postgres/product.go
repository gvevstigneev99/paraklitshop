package postgres

import (
	"paraklitshop/internal/model"
	"paraklitshop/internal/repository"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) repository.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts() ([]model.Product, error) {
	const query = `
		SELECT id, title, description, price, seller_id
		FROM products
		ORDER BY id
	`
	var products []model.Product
	err := r.db.Select(&products, query)
	if err != nil {
		return nil, err
	}
	return products, nil
}
