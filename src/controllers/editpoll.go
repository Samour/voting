package controllers

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/polls"
)

var renderer = Must(CreateRenderer("../resources/pages/*.html"))

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

	err = renderer.Render(w, "edit_poll.html", poll)
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

	poll, err := polls.UpdatePollDetails(pollId, polls.PollDetails{
		Name:    name,
		Options: options,
	})
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
	err = renderer.Render(w, "view_poll.html", poll)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

type optionsForm struct {
	Options []string
}

func HandlePatchPoll(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	options := r.PostForm["Options[]"]
	add := r.PostForm.Has("Add")
	remove := -1
	if r.PostForm.Has("Remove") {
		remove, err = strconv.Atoi(r.PostForm.Get("Remove"))
		if err != nil {
			errorPage(w, err.Error(), http.StatusBadRequest)
		}
	}

	options = polls.PatchPollOptions(options, polls.PollOptionsUpdate{
		Add:    add,
		Remove: remove,
	})
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "edit_poll_options.html", optionsForm{
		Options: options,
	})
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePollStatusChange(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		errorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := r.PostForm.Get("Status")
	poll, err := polls.UpdateStatus(pollId, status)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "view_poll_navigation.html", poll)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
