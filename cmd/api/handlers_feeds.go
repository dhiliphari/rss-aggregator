package main

import (
	"encoding/json"
	"net/http"

	"github.com/dhiliphari/rss-aggregator/internal/models"
)

func (app *application) createFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		panic("handler missing auth middleware")
	}

	var params struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	feed, err := app.feeds.Create(r.Context(), params.Name, params.URL, user.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, feed)
}
