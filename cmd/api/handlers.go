package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/index.html")
	if err != nil {
		app.serverError(w, err)
		return
	}

	ts.Execute(w, nil)
}

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, struct{}{})
}
