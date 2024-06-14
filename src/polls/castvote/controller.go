package castvote

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/poll_vote.html"))

type PollVoteForm struct {
	Poll  *model.Poll
	Voted int
}

func ServeVotePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	poll, err := repository.GetPollItem(pollId)
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

	poll, err := castVote(pollId, option)
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
