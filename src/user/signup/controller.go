package signup

import (
	"net/http"

	"github.com/Samour/voting/render"
)

var renderer = render.Must(render.CreateRenderer("pages/signup.html"))

func ServeSignUp(w http.ResponseWriter, r *http.Request) {
	renderer.UsingTemplate(w, "signup.html").Render(render.HttpResponse{}, nil)
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		render.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("Username")
	password := r.PostForm.Get("Password")

	redirect, res, err := createAccount(username, password)
	if redirect != nil {
		http.Redirect(w, r, *redirect, http.StatusFound)
	} else {
		renderer.UsingTemplate(w, "signup.html").Render(res, err)
	}
}
