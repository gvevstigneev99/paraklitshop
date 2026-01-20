package model

import "time"

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalPrice float64   `json:"total"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
