package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/dhiliphari/rss-aggregator/internal/models"
)

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := getApiKey(r)
		if err != nil {
			app.clientError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := app.users.Get(r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				app.clientError(w, http.StatusNotFound, "user not found")
			} else {
				app.serverError(w, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
