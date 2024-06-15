package viewpoll

import (
	"net/http"
	"strings"

	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/view_poll/*.html"))

func ServeViewPoll(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	isHxRequest := strings.ToLower(r.Header.Get("HX-Request")) == "true"

	poll, err := getPoll(pollId)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if poll == nil {
		render.ErrorPage(w, "Poll not found", http.StatusNotFound)
		return
	}

	poll.RenderFullPage = !isHxRequest
	err = renderer.Render(w, "index.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePollStatusChange(w http.ResponseWriter, r *http.Request) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := r.PostForm.Get("Status")
	poll, err := updateStatus(pollId, status)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	poll.RenderFullPage = false
	err = renderer.Render(w, "index.html", poll)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
