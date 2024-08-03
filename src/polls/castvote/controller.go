package castvote

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/poll_vote/*.html"))

func ServeVotePoll(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetSession(r)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pollId := r.PathValue("id")
	err = r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderer.UsingTemplate(w, "index.html").Render(getPollVoteForm(session, pollId))
}

func HandleCastFptpVote(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetSession(r)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pollId := r.PathValue("id")
	err = r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	option := -1
	option, err = strconv.Atoi(r.PostForm.Get("Option"))
	if err != nil {
		option = -1
	}

	renderer.UsingTemplate(w, "vote_form.html").Render(castFptpVote(session, pollId, option))
}

func HandlePatchRankedChoice(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	selected, err := extractSelectedArray(r.PostForm)
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

	renderer.UsingTemplate(w, "rankedchoice.html").Render(
		updateRankedChoiceOption(pollId, selected, rankedChoiceUpdate{
			Select:   newSelection,
			Unselect: remove,
		}))
}

func HandleCastRankedChoiceVote(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetSession(r)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pollId := r.PathValue("id")

	err = r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	ranked, err := extractSelectedArray(r.PostForm)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderer.UsingTemplate(w, "vote_form.html").Render(castRankedChoiceVote(session, pollId, ranked))
}

func extractSelectedArray(v url.Values) ([]int, error) {
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
