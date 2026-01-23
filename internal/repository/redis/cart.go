package redis

import (
	"paraklitshop/internal/repository"
	"sync"
)

type CartRepository struct {
	mu    sync.Mutex
	carts map[int]map[int]int // userID -> productID -> quantity
}

func NewCartRepository(addr, password string, db int) (repository.CartRepository, error) {
	// TODO: Implement real Redis connection when redis library is added
	// For now using in-memory storage
	return &CartRepository{
		carts: make(map[int]map[int]int),
	}, nil
}

func (r *CartRepository) Add(userID, productID, qty int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.carts[userID] == nil {
		r.carts[userID] = make(map[int]int)
	}
	r.carts[userID][productID] += qty
	return nil
}

func (r *CartRepository) Get(userID int) (map[int]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart := make(map[int]int)
	if r.carts[userID] != nil {
		for productID, quantity := range r.carts[userID] {
			cart[productID] = quantity
		}
	}
	return cart, nil
}

func (r *CartRepository) Clear(userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.carts, userID)
	return nil
}

func (r *CartRepository) Remove(userID, productID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.carts[userID] != nil {
		delete(r.carts[userID], productID)
	}
	return nil
}
