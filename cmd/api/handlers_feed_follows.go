package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"uuid"

	"github.com/dhiliphari/rss-aggregator/internal/models"
	"github.com/dhiliphari/rss-aggregator/internal/validator"
)

func (app *application) followFeed(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)

	var params struct {
		FeedId string `json:"feed_id"`
		validator.Validator
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.CheckField(validator.NotBlank(params.FeedId), "feed_id", "field is required")

	feedId, err := uuid.Parse(params.FeedId)
	if err != nil {
		params.AddFieldError("feed_id", "must be a valid uuid")
	}

	if !params.Valid() {
		app.clientError(w, http.StatusBadRequest, params.FieldErrors)
		return
	}

	feedFollow, err := app.feedFollows.Create(r.Context(), user.Id, feedId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, feedFollow)
}

func (app *application) feedsFollowedByUser(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)

	feedFollows, err := app.feedFollows.ListByUser(r.Context(), user.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, feedFollows)
}

func (app *application) unfollowFeed(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)

	feedFollowId, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid feed follow id")
		return
	}

	err = app.feedFollows.Delete(r.Context(), feedFollowId, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			app.clientError(w, http.StatusNotFound, "feed follow not found")
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusNoContent, struct{}{})
}
