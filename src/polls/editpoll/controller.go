package editpoll

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/edit_poll/*.html"))
var viewPollRenderer = render.Must(render.CreateRenderer("pages/view_poll/*.html"))

func ServeEditPoll(w http.ResponseWriter, r *http.Request) {
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

	// TODO handle poll that is not in draft

	err = renderer.Render(w, "index.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeSavePoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("Name")
	options := r.PostForm["Options[]"]

	poll, err := updatePollDetails(pollId, pollDetails{
		Name:    name,
		Options: options,
	})
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		render.ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	err = viewPollRenderer.Render(w, "index.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

type optionsForm struct {
	Options []string
}

func HandlePatchPoll(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	options := r.PostForm["Options[]"]
	add := r.PostForm.Has("Add")
	remove := -1
	if r.PostForm.Has("Remove") {
		remove, err = strconv.Atoi(r.PostForm.Get("Remove"))
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		}
	}

	options = patchPollOptions(options, pollOptionsUpdate{
		Add:    add,
		Remove: remove,
	})
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "poll_options.html", optionsForm{
		Options: options,
	})
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
