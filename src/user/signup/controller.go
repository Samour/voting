package signup

import (
	"net/http"

	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/signup.html"))

func ServeSignUp(w http.ResponseWriter, r *http.Request) {
	renderer.UsingTemplate(w, "signup.html").Render(prepareSignUp())
}

func prepareSignUp() (render.HttpResponse, error) {
	return render.HttpResponse{}, nil
}
