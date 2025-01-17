package rest

import (
	"context"
	"lib/internal/domain"
	"net/http"

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
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	ParseToken(ctx context.Context, accessToken string) (int64, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
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
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)
	}

	books := r.PathPrefix("/books").Subrouter()
	{
		books.Use(h.authMiddleware)

		books.HandleFunc("/", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("/", h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods(http.MethodPut)
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.getBookByID).Methods(http.MethodGet)
	}

	return r
}
