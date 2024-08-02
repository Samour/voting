package login

import (
	"net/http"

	"github.com/Samour/voting/auth"
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

	loginSuccess, page, err := logIn(username, password)
	if len(loginSuccess.SessionId) > 0 {
		auth.WriteSessionCookie(w, loginSuccess.SessionId)
		http.Redirect(w, r, loginSuccess.Redirect, http.StatusFound)
	} else {
		renderer.UsingTemplate(w, "login.html").Render(page, err)
	}
}

func ServeLogOut(w http.ResponseWriter, r *http.Request, session auth.Session) {
	auth.RemoveSession(session.SessionId)
	auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
