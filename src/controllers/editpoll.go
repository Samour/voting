package controllers

import (
	"net/http"
)

func ServeEditPoll(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "edit_poll.html", nil)
	if err != nil {
		ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
