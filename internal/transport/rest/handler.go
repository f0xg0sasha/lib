package rest

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lib/internal/domain"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Books interface {
	Create(ctx context.Context, book domain.Book) error
	Update(ctx context.Context, id int64, inp domain.UpdateBook) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetByID(ctx context.Context, id int64) (domain.Book, error)
}

type Handler struct {
	booksService Books
}

func NewBooksHandler(books Books) *Handler {
	return &Handler{
		booksService: books,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

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

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book := domain.Book{}
	err = json.Unmarshal(reqBytes, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.booksService.Create(context.TODO(), book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"error":   err,
		}).Error()

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.booksService.GetAll(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getAllBooks",
			"error":   err,
		}).Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book := domain.UpdateBook{}
	err = json.Unmarshal(reqBytes, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = h.booksService.Update(context.TODO(), id, book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"error":   err,
		}).Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = h.booksService.Delete(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"error":   err,
		}).Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	books, err := h.booksService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"error":   err,
		}).Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(books)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"error":   err,
		}).Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
