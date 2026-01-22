package repository

import (
	"context"
	"database/sql"
	"errors"
	"paraklitshop/internal/model"

	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type UserPostgresRepository struct {
	db *sqlx.DB
}

func NewUserPostgresRepository(db *sqlx.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgresRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password_hash, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role, user.CreatedAt).Scan(&user.ID)
	return err
}
