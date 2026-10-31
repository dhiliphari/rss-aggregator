package main

import (
	"encoding/json"
	"net/http"

	"github.com/dhiliphari/rss-aggregator/internal/models"
	"github.com/dhiliphari/rss-aggregator/internal/validator"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Name string `json:"name"`
		validator.Validator
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.CheckField(validator.NotBlank(params.Name), "name", "field is required")

	if !params.Valid() {
		app.clientError(w, http.StatusBadRequest, params.FieldErrors)
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
