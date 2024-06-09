package controllers

import (
	"fmt"
	"net/http"

	"github.com/Samour/voting/polls"
)

func ServeNewPoll(w http.ResponseWriter, r *http.Request) {
	id := polls.CreatePoll()

	redirect := fmt.Sprintf("/polls/%s/edit", id)
	http.Redirect(w, r, redirect, http.StatusFound)
}
