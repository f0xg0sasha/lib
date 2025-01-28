package psql

import (
	"context"
	"database/sql"
	"lib/internal/domain"

	"github.com/sirupsen/logrus"
)

type Token struct {
	db *sql.DB
}

func NewToken(db *sql.DB) *Token {
	return &Token{
		db: db,
	}
}

func (t *Token) Create(ctx context.Context, token domain.RefreshSession) error {
	_, err := t.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt,
	)

	return err
}

func (t *Token) Get(ctx context.Context, token string) (domain.RefreshSession, error) {
	var session domain.RefreshSession

	err := t.db.QueryRow("SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		logrus.Info("refresh token no rows")
		return session, nil
	}
	if err != nil {
		logrus.Info("refresh token error")
		return session, err
	}

	_, err = t.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token=$1", token)

	return session, err
}
