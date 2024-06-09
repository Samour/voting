package controllers

import (
	"net/http"

	"github.com/Samour/voting/polls"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	polls, err := polls.FetchAllPolls()
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderTemplate(w, "home.html", polls)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
