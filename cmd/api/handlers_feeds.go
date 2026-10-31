package main

import (
	"encoding/json"
	"net/http"

	"github.com/dhiliphari/rss-aggregator/internal/models"
	"github.com/dhiliphari/rss-aggregator/internal/validator"
)

func (app *application) createFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		panic("handler missing auth middleware")
	}

	var params struct {
		Name string `json:"name"`
		URL  string `json:"url"`
		validator.Validator
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.CheckField(validator.NotBlank(params.Name), "name", "field is required")
	params.CheckField(validator.NotBlank(params.URL), "url", "field is required")
	params.CheckField(validator.IsValidURL(params.URL), "url", "must be a valid URL")

	if !params.Valid() {
		app.clientError(w, http.StatusBadRequest, params.FieldErrors)
		return
	}

	feed, err := app.feeds.Create(r.Context(), params.Name, params.URL, user.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, feed)
}

func (app *application) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.feeds.List(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, feeds)
}
