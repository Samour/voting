package home

import (
	"net/http"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/polls/getpolls"
	"github.com/Samour/voting/render"
	"github.com/Samour/voting/site"
)

var renderer = render.Must(render.CreateRenderer("pages/home.html"))

func ServeHome(w http.ResponseWriter, r *http.Request, s auth.Session) {
	renderer.UsingTemplate(w, "home.html").Render(prepareHome(s))
}

func prepareHome(s auth.Session) (render.HttpResponse, error) {
	polls, err := getpolls.FetchAllPolls()
	if err != nil {
		return render.HttpResponse{}, err
	}

	return render.HttpResponse{
		Model: homeModel{
			SiteModel: site.BuildSiteModel(s),
			Polls:     polls,
		},
	}, nil
}
