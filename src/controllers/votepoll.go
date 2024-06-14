package controllers

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/polls"
	"github.com/Samour/voting/render"
)

type PollVoteForm struct {
	Poll  *polls.Poll
	Voted int
}

func ServeVotePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	poll, err := polls.FetchPoll(pollId)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		render.ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	f := PollVoteForm{
		Poll:  poll,
		Voted: -1,
	}
	err = renderer.Render(w, "poll_vote.html", f)
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

	poll, err := polls.CastVote(pollId, option)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		render.ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	f := PollVoteForm{
		Poll:  poll,
		Voted: option,
	}
	err = renderer.Render(w, "poll_vote.html", f)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
