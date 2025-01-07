package service

import (
	"context"
	"lib/internal/domain"
	"time"
)

type BooksRepository interface {
	Create(ctx context.Context, book domain.Book) error
	Update(ctx context.Context, id int64, inp domain.UpdateBook) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetByID(ctx context.Context, id int64) (domain.Book, error)
}

type Books struct {
	repo BooksRepository
}

func NewBooks(repo BooksRepository) *Books {
	return &Books{
		repo: repo,
	}
}

func (b *Books) Create(ctx context.Context, book domain.Book) error {
	if book.Publisher.IsZero() {
		book.Publisher = time.Now()
	}

	return b.repo.Create(ctx, book)
}

func (b *Books) Update(ctx context.Context, id int64, inp domain.UpdateBook) error {
	return b.repo.Update(ctx, id, inp)
}

func (b *Books) Delete(ctx context.Context, id int64) error {
	return b.repo.Delete(ctx, id)
}

func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
	return b.repo.GetAll(ctx)
}

func (b *Books) GetByID(ctx context.Context, id int64) (domain.Book, error) {
	return b.repo.GetByID(ctx, id)
}
