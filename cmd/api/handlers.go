package main

import "net/http"

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, struct{}{})
}
