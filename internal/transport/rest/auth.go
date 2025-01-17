package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lib/internal/domain"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.SignUpInput
	if err := json.Unmarshal(reqBytes, &inp); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.usersService.SignUp(r.Context(), inp); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.SignInInput
	if err := json.Unmarshal(reqBytes, &inp); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.usersService.SignIn(r.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(w, err)
			return
		}

		logError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func handleNotFoundError(w http.ResponseWriter, err error) {
	response, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		logError("refresh", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.Infof("%s", cookie.Value)

	accsesToken, refreshToken, err := h.usersService.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		logError("refresh", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accsesToken,
	})

	if err != nil {
		logError("refresh", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token='%s'; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}
