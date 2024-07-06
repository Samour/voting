package home

import (
	"net/http"

	"github.com/Samour/voting/polls/getpolls"
	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/home.html"))

func ServeHome(w http.ResponseWriter, r *http.Request) {
	renderer.UsingTemplate(w, "home.html").Render(prepareHome())
}

func prepareHome() (render.HttpResponse, error) {
	polls, err := getpolls.FetchAllPolls()
	if err != nil {
		return render.HttpResponse{}, err
	}

	return render.HttpResponse{
		Model: homeModel{
			Polls: polls,
		},
	}, nil
}
