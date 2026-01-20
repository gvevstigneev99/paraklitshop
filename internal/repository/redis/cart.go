package redis

import "sync"

type CartRepository struct {
	mu    sync.Mutex
	carts map[int]map[int]int // userID -> productID -> quantity
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		carts: make(map[int]map[int]int),
	}
}

func (r *CartRepository) AddItem(userID, productID, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.carts[userID] == nil {
		r.carts[userID] = make(map[int]int)
	}
	r.carts[userID][productID] += quantity
	return nil
}

func (r *CartRepository) GetCart(userID int) (map[int]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart := make(map[int]int)
	for productID, quantity := range r.carts[userID] {
		cart[productID] = quantity
	}
	return cart, nil
}

func (r *CartRepository) ClearCart(userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.carts, userID)
	return nil
}
func (r *CartRepository) RemoveItem(userID, productID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.carts[userID] != nil {
		delete(r.carts[userID], productID)
	}
	return nil
}
