package repository

import (
	"context"
	"errors"
	"paraklitshop/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository описывает работу с пользователями
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int, error)
}

// ProductRepository описывает работу с товарами
type ProductRepository interface {
	GetAllProducts() ([]model.Product, error)
}

// OrderRepository описывает работу с заказами
type OrderRepository interface {
	Create(order *model.Order) error
}

// CartRepository описывает работу с корзиной (Redis)
type CartRepository interface {
	Add(userID, productID, qty int) error
	Remove(userID, productID int) error
	Clear(userID int) error
	Get(userID int) (map[int]int, error) // productID -> quantity
}
