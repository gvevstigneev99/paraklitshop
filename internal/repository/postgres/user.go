package postgres

import (
	"context"
	"database/sql"
	"errors"
	"paraklitshop/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	const query = `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`
	var u model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int, error) {
	const query = `
		INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
