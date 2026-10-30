package main

import (
	"encoding/json"
	"net/http"

	"github.com/dhiliphari/rss-aggregator/internal/models"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := app.users.Create(r.Context(), params.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, user)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		panic("handler missing auth middleware")
	}
	app.writeJSON(w, http.StatusOK, user)
}
