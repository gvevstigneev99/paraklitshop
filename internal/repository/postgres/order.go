package postgres

import "paraklitshop/internal/model"

type OrderRepository struct {
	orders []model.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: []model.Order{},
	}
}

func (r *OrderRepository) CreateOrder(order model.Order) error {
	r.orders = append(r.orders, order)
	return nil
}
