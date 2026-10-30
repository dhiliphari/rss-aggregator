package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.healthz)

	mux.HandleFunc("POST /users", app.createUser)
	mux.Handle("GET /users", app.requireAuthentication(http.HandlerFunc(app.getUser)))

	return mux
}
