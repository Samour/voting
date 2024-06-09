package controllers

import (
	"net/http"

	"github.com/Samour/voting/polls"
)

func ServeViewPoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	poll, err := polls.FetchPoll(pollId)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		errorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	err = renderTemplate(w, "view_poll.html", poll)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
