package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Samour/voting/polls"
)

func ServeEditPoll(w http.ResponseWriter, r *http.Request) {
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

	// TODO handle poll that is not in draft

	err = renderTemplate(w, "edit_poll.html", poll)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeSavePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		errorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("Name")
	options := r.PostForm["Options[]"]

	poll, err := polls.UpdatePollDetails(pollId, name, options)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		errorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	// TODO render "view" screen
	// For now, just redirect back to edit screen
	redirect := fmt.Sprintf("/polls/%s/edit", pollId)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func HandlePatchPoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		errorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("Name")
	options := r.PostForm["Options[]"]
	add := r.PostForm.Has("Add")
	remove := -1
	if r.PostForm.Has("Remove") {
		remove, err = strconv.Atoi(r.PostForm.Get("Remove"))
		if err != nil {
			errorPage(w, err.Error(), http.StatusBadRequest)
		}
	}

	poll, err := polls.PatchPollOptions(pollId, polls.PollOptionsUpdate{
		Details: polls.PollDetails{
			Name:    name,
			Options: options,
		},
		Add:    add,
		Remove: remove,
	})
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderTemplate(w, "edit_poll_options.html", poll)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
