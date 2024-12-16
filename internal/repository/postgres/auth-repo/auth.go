package authrepo

import (
	"context"
	"kn-assignment/internal/core/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repository) CreateUser(ctx context.Context, user domain.CreateUserRequest) error {
	query := `INSERT INTO users (username, password, role, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.dbPool.Exec(ctx, query, user.Username, user.Password, user.Role)
	return err
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	var user domain.User
	err := pgxscan.Get(ctx, r.dbPool, &user, query, username)
	if pgxscan.NotFound(err) {
		return nil, nil
	}
	return &user, err
}

func (r *repository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user domain.User
	err := pgxscan.Get(ctx, r.dbPool, &user, query, userID)
	return &user, err
}

func (r *repository) UpdateUser(ctx context.Context, user domain.User) error {
	query := `UPDATE users SET username = $1, password = $2, role = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.dbPool.Exec(ctx, query, user.Username, user.Password, user.Role, user.ID)
	return err
}
