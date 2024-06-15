package castvote

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/poll_vote/*.html"))

func ServeVotePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	var voteId *string = nil
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Form.Has("id") {
		v := r.Form.Get("id")
		voteId = &v
	}

	poll, err := getPollVoteForm(pollId, voteId)
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

func ServeCastFptpVote(w http.ResponseWriter, r *http.Request) {
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

	voteId, err := castVote(pollId, option)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirect := fmt.Sprintf("/polls/%s/vote/?id=%s", pollId, *voteId)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func HandlePatchRankedChoice(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	selected, err := extractSelectedArray(&r.PostForm)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	newSelection := -1
	if r.PostForm.Has("Select") {
		newSelection, err = strconv.Atoi(r.PostForm.Get("Select"))
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	remove := -1
	if r.PostForm.Has("Remove") {
		remove, err = strconv.Atoi(r.PostForm.Get("Remove"))
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	poll, err := updateRankedChoiceOption(pollId, selected, rankedChoiceUpdate{
		Select:   newSelection,
		Unselect: remove,
	})
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "rankedchoice.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func extractSelectedArray(v *url.Values) ([]int, error) {
	r := make([]int, 0)
	i := 0
	for {
		key := fmt.Sprintf("Selected[%d]", i)
		if !v.Has(key) {
			return r, nil
		}

		value, err := strconv.Atoi(v.Get(key))
		if err != nil {
			return nil, err
		}
		r = append(r, value)
		i++
	}
}
