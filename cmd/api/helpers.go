package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dhiliphari/rss-aggregator/internal/models"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (app *application) clientError(w http.ResponseWriter, status int, msg any) {
	payload := struct {
		Error any `json:"error"`
	}{
		Error: msg,
	}
	app.writeJSON(w, status, payload)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func getApiKey(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("no auth header")
	}

	split := strings.Split(auth, " ")
	if len(split) != 2 || split[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	return split[1], nil
}

func getUser(r *http.Request) *models.User {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		panic("handler missing auth middleware")
	}
	return user
}
