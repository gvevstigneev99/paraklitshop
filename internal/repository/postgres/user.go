package postgres

import (
	"context"
	"database/sql"
	"paraklitshop/internal/model"
	"paraklitshop/internal/repository"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	var u model.User
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int, error) {
	const query = `
		INSERT INTO users (email, password_hash, role, created_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
