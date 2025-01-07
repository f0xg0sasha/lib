package domain

import (
	"errors"
	"time"
)

var (
	ErrBookNotFound = errors.New("Book not found")
)

type Book struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	Publisher time.Time `json:"publisher"`
	Rating    int       `json:"rating"`
}

type UpdateBook struct {
	Name      *string    `json:"name"`
	Author    *string    `json:"author"`
	Publisher *time.Time `json:"publisher"`
	Rating    *int       `json:"rating"`
}
