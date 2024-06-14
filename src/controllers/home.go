package controllers

import (
	"net/http"

	"github.com/Samour/voting/polls"
	"github.com/Samour/voting/render"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	polls, err := polls.FetchAllPolls()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "home.html", polls)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
