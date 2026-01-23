package postgres

import (
	"paraklitshop/internal/model"
	"paraklitshop/internal/repository"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) repository.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	const query = `
		INSERT INTO orders (user_id, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		order.UserID,
		order.TotalPrice,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)
	return err
}
