package main

import (
	"net/http"
	"strconv"
)

func (app *application) listUserPosts(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)

	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		count = 10
	}

	posts, err := app.posts.ListByUser(r.Context(), user.Id, count)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, posts)
}
