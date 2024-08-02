package viewpoll

import (
	"net/http"
	"strings"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/view_poll/*.html"))

func ServeViewPoll(w http.ResponseWriter, r *http.Request, s auth.Session) {
	pollId := r.PathValue("id")
	isHxRequest := strings.ToLower(r.Header.Get("HX-Request")) == "true"

	renderer.UsingTemplate(w, "index.html").Render(getPoll(pollId, !isHxRequest))
}

func HandlePollStatusChange(w http.ResponseWriter, r *http.Request, s auth.Session) {
	pollId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := r.PostForm.Get("Status")

	renderer.UsingTemplate(w, "index.html").Render(updateStatus(pollId, status))
}
