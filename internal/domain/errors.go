package domain

import "errors"

var (
	ErrBookNotFound        = errors.New("book not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
