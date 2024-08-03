package editpoll

import (
	"net/http"
	"strconv"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/edit_poll/*.html"))
var viewPollRenderer = render.Must(render.CreateRenderer("pages/view_poll/*.html"))

func ServeEditPoll(w http.ResponseWriter, r *http.Request, s auth.Session) {
	pollId := r.PathValue("id")

	renderer.UsingTemplate(w, "index.html").Render(getPoll(s, pollId))
}

func ServeSavePoll(w http.ResponseWriter, r *http.Request, s auth.Session) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("Name")
	aggregationType := r.PostForm.Get("AggregationType")
	options := r.PostForm["Options[]"]

	viewPollRenderer.UsingTemplate(w, "index.html").Render(
		updatePollDetails(s, pollId, pollDetails{
			Name:            name,
			AggregationType: aggregationType,
			Options:         options,
		}))
}

func HandlePatchPoll(w http.ResponseWriter, r *http.Request, s auth.Session) {
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
			return
		}
	}

	renderer.UsingTemplate(w, "poll_options.html").Render(
		patchPollOptions(options, pollOptionsUpdate{
			Add:    add,
			Remove: remove,
		}), nil)
}
