package main

import (
	"encoding/json"
	"net/http"
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

	user, err := app.users.CreateUser(r.Context(), params.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, user)
}
