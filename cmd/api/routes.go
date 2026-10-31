package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// get server health status
	mux.HandleFunc("GET /healthz", app.healthz)

	// user: create
	mux.Handle("POST /users", http.HandlerFunc(app.createUser))
	// user: get by api_key
	mux.Handle("GET /me", app.requireAuthentication(http.HandlerFunc(app.getUser)))

	// feed: list all
	mux.Handle("GET /feeds", http.HandlerFunc(app.listFeeds))
	// feed: create by user
	mux.Handle("POST /me/feeds", app.requireAuthentication(http.HandlerFunc(app.createFeed)))
	// feed: get by user
	mux.Handle("GET /me/feeds", app.requireAuthentication(http.HandlerFunc(app.listUserFeeds)))

	return mux
}
