package login

import (
	"net/http"

	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/login.html"))

func ServeLogIn(w http.ResponseWriter, r *http.Request) {
	renderer.UsingTemplate(w, "login.html").Render(render.HttpResponse{}, nil)
}

func HandleLogIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("Username")
	password := r.PostForm.Get("Password")

	redirect, res, err := logIn(username, password)
	if redirect != nil {
		http.Redirect(w, r, *redirect, http.StatusFound)
	} else {
		renderer.UsingTemplate(w, "login.html").Render(res, err)
	}
}
