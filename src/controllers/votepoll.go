package controllers

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/polls"
)

type PollVoteForm struct {
	Poll  *polls.Poll
	Voted int
}

func ServeVotePoll(w http.ResponseWriter, r *http.Request) {
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

	f := PollVoteForm{
		Poll:  poll,
		Voted: -1,
	}
	err = renderTemplate(w, "poll_vote.html", f)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeCastVote(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		errorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	option, err := strconv.Atoi(r.PostForm.Get("Option"))
	if err != nil {
		errorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	poll, err := polls.CastVote(pollId, option)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		errorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	f := PollVoteForm{
		Poll:  poll,
		Voted: option,
	}
	err = renderTemplate(w, "poll_vote.html", f)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
