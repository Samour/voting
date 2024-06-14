package controllers

import (
	"fmt"
	"net/http"

	"github.com/Samour/voting/polls"
	"github.com/Samour/voting/render"
)

func ServeNewPoll(w http.ResponseWriter, r *http.Request) {
	id, err := polls.CreatePoll()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirect := fmt.Sprintf("/polls/%s/edit", *id)
	http.Redirect(w, r, redirect, http.StatusFound)
}
