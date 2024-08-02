package signup

import (
	"net/http"

	"github.com/Samour/voting/auth"
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

	signUpSuccess, page, err := createAccount(username, password)
	if len(signUpSuccess.SessionId) > 0 {
		auth.WriteSessionCookie(w, signUpSuccess.SessionId)
		http.Redirect(w, r, signUpSuccess.Redirect, http.StatusFound)
	} else {
		renderer.UsingTemplate(w, "signup.html").Render(page, err)
	}
}
