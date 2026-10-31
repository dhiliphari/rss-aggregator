package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// get server health status
	mux.HandleFunc("GET /healthz", app.healthz)

	// users: create
	mux.Handle("POST /users", http.HandlerFunc(app.createUser))
	// users: get by api_key
	mux.Handle("GET /me", app.requireAuthentication(http.HandlerFunc(app.getUser)))

	// feeds: list all
	mux.Handle("GET /feeds", http.HandlerFunc(app.listFeeds))
	// feeds: create by user
	mux.Handle("POST /me/feeds", app.requireAuthentication(http.HandlerFunc(app.createFeed)))
	// feeds: get by user
	mux.Handle("GET /me/feeds", app.requireAuthentication(http.HandlerFunc(app.listUserFeeds)))

	// feed_follows: feeds followed by user
	mux.Handle("GET /me/feeds/following", app.requireAuthentication(http.HandlerFunc(app.feedsFollowedByUser)))
	// feed_follows: user follow feed
	mux.Handle("POST /me/feeds/following", app.requireAuthentication(http.HandlerFunc(app.followFeed)))
	// feed_follows: user follow feed
	mux.Handle("DELETE /me/feeds/following/{feedFollowID}", app.requireAuthentication(http.HandlerFunc(app.unfollowFeed)))

	// posts: feeds followed by user
	mux.Handle("GET /me/posts", app.requireAuthentication(http.HandlerFunc(app.listUserPosts)))

	return mux
}
