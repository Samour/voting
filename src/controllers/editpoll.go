package controllers

import (
	"net/http"

	"github.com/Samour/voting/polls"
)

func ServeEditPoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	poll, err := polls.FetchPoll(pollId)
	if err != nil {
		ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	err = Templates.ExecuteTemplate(w, "edit_poll.html", poll)
	if err != nil {
		ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
