package castvote

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/poll_vote/*.html"))

func ServeVotePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	poll, err := getPoll(pollId)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		render.ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	err = renderer.Render(w, "index.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeCastVote(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	option, err := strconv.Atoi(r.PostForm.Get("Option"))
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	poll, err := castVote(pollId, option)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		render.ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	err = renderer.Render(w, "index.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePatchRankedChoice(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	selected, err := strconv.Atoi(r.PostForm.Get("Select"))
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	poll, err := selectRankedChoiceOption(pollId, selected)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "rankedchoice.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
