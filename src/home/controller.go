package home

import (
	"net/http"

	"github.com/Samour/voting/polls/getpolls"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("../resources/pages/*.html"))

func ServeHome(w http.ResponseWriter, r *http.Request) {
	polls, err := getpolls.FetchAllPolls()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderer.Render(w, "home.html", polls)
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
