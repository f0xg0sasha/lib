package psql

import (
	"context"
	"database/sql"
	"lib/internal/domain"
)

type User struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *User {
	return &User{db: db}
}

func (r *User) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)
	return err
}

func (r *User) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email, password, registered_at FROM users WHERE email=$1 AND password=$2",
		email, password).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RegisteredAt)

	return user, err
}
