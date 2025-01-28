package service

import (
	"context"
	"lib/internal/domain"
	"time"

	"github.com/f0xg0sasha/audit_logger/pkg/domain/audit"
)

type BooksRepository interface {
	Create(ctx context.Context, book domain.Book) (int64, error)
	Update(ctx context.Context, id int64, inp domain.UpdateBook) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetByID(ctx context.Context, id int64) (domain.Book, error)
}

type Books struct {
	repo        BooksRepository
	auditClient AuditClient
}

func NewBooks(repo BooksRepository, auditClient AuditClient) *Books {
	return &Books{
		repo:        repo,
		auditClient: auditClient,
	}
}

func (b *Books) Create(ctx context.Context, book domain.Book) error {
	if book.Publisher.IsZero() {
		book.Publisher = time.Now()
	}

	id, err := b.repo.Create(ctx, book)
	if err != nil {
		return err
	}

	err = b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Entity:    audit.ENTITY_BOOK,
		Action:    audit.ACTION_CREATE,
		EntityID:  id,
		Timestamp: time.Now(),
	})

	return err
}

func (b *Books) Update(ctx context.Context, id int64, inp domain.UpdateBook) error {
	err := b.repo.Update(ctx, id, inp)
	if err != nil {
		return err
	}

	err = b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Entity:    audit.ENTITY_BOOK,
		Action:    audit.ACTION_UPDATE,
		EntityID:  id,
		Timestamp: time.Now(),
	})

	return err
}

func (b *Books) Delete(ctx context.Context, id int64) error {
	err := b.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Entity:    audit.ENTITY_BOOK,
		Action:    audit.ACTION_DELETE,
		EntityID:  id,
		Timestamp: time.Now(),
	})

	return err
}

func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
	books, err := b.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	err = b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Entity:    audit.ENTITY_BOOK,
		Action:    audit.ACTION_GET,
		EntityID:  0,
		Timestamp: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b *Books) GetByID(ctx context.Context, id int64) (domain.Book, error) {
	book, err := b.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}

	err = b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Entity:    audit.ENTITY_BOOK,
		Action:    audit.ACTION_GET,
		EntityID:  id,
		Timestamp: time.Now(),
	})

	if err != nil {
		return domain.Book{}, err
	}

	return book, nil
}
