package rest

import (
	"context"
	"lib/internal/domain"

	"github.com/gorilla/mux"
)

type Books interface {
	Create(ctx context.Context, book domain.Book) error
	Update(ctx context.Context, id int64, inp domain.UpdateBook) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetByID(ctx context.Context, id int64) (domain.Book, error)
}

type User interface {
	SignIn(ctx context.Context, inp domain.SignInInput) (string, error)
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	ParseToken(ctx context.Context, token string) (int64, error)
}

type Handler struct {
	booksService Books
	usersService User
}

func NewHandler(books Books, users User) *Handler {
	return &Handler{
		booksService: books,
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods("POST")
		auth.HandleFunc("/sign-in", h.signIn).Methods("POST")
	}

	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("/", h.createBook).Methods("POST")
		books.HandleFunc("/", h.getAllBooks).Methods("GET")
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods("PUT")
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods("DELETE")
		books.HandleFunc("/{id:[0-9]+}", h.getBookByID).Methods("GET")
	}

	return r
}
