package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.healthz)

	mux.HandleFunc("POST /users", app.createUser)

	return mux
}
