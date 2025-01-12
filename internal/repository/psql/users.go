package psql

import (
	"context"
	"database/sql"
	"lib/internal/domain"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, user domain.User) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)
	return err
}

func (u *User) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	return domain.User{}, nil
}
