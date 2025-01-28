package psql

import (
	"context"
	"database/sql"
	"fmt"
	"lib/internal/domain"
	"strings"
)

type Books struct {
	db *sql.DB
}

func NewBooks(db *sql.DB) *Books {
	return &Books{
		db: db,
	}
}

func (b *Books) Create(ctx context.Context, book domain.Book) (int64, error) {
	var id int64
	err := b.db.QueryRow("INSERT INTO books (name, author, publisher, rating) VALUES ($1, $2, $3, $4) RETURNING id",
		book.Name, book.Author, book.Publisher, book.Rating).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := b.db.QueryContext(ctx, "SELECT id, name, author, publisher, rating FROM books")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := make([]domain.Book, 0)
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Publisher, &book.Rating); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func (b *Books) GetByID(ctx context.Context, id int64) (domain.Book, error) {
	row := b.db.QueryRowContext(ctx, "SELECT id, name, author, publisher, rating FROM books WHERE id = $1", id)

	var book domain.Book
	if err := row.Scan(&book.ID, &book.Name, &book.Author, &book.Publisher, &book.Rating); err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func (b *Books) Update(ctx context.Context, id int64, inp domain.UpdateBook) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsID := 1

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argsID))
		args = append(args, *inp.Name)
		argsID++
	}

	if inp.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author = $%d", argsID))
		args = append(args, *inp.Author)
		argsID++
	}

	if inp.Publisher != nil {
		setValues = append(setValues, fmt.Sprintf("publisher = $%d", argsID))
		args = append(args, *inp.Publisher)
		argsID++
	}

	if inp.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating = $%d", argsID))
		args = append(args, *inp.Rating)
		argsID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE books SET %s WHERE id = $%d", setQuery, argsID)
	args = append(args, id)

	_, err := b.db.ExecContext(ctx, query, args...)
	return err
}

func (b *Books) Delete(ctx context.Context, id int64) error {
	_, err := b.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	return err
}
