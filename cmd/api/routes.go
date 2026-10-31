package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.healthz)

	mux.Handle("POST /users", http.HandlerFunc(app.createUser))
	mux.Handle("GET /users", app.requireAuthentication(http.HandlerFunc(app.getUser)))

	mux.Handle("GET /feeds", http.HandlerFunc(app.listFeeds))
	mux.Handle("POST /feeds", app.requireAuthentication(http.HandlerFunc(app.createFeed)))

	return mux
}
