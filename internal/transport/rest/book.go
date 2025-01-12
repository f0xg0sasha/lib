package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"lib/internal/domain"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

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
		logError("CreateBook", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.booksService.GetAll(context.TODO())
	if err != nil {
		logError("GetAllBooks", err)
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
		logError("UpdateBook", err)

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
		logError("DeleteBook", err)
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

		logError("GetBookByID", err)
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
