package login

import (
	"net/http"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/middleware"
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
	redirect := middleware.GetAuthRedirect(r)

	loginSuccess, page, err := logIn(username, password, redirect)
	if len(loginSuccess.SessionId) > 0 {
		auth.WriteSessionCookie(w, loginSuccess.SessionId)
		http.Redirect(w, r, loginSuccess.Redirect, http.StatusFound)
	} else {
		renderer.UsingTemplate(w, "login.html").Render(page, err)
	}
}

func ServeLogOut(w http.ResponseWriter, r *http.Request, s auth.Session) {
	auth.RemoveSession(s.SessionId)
	auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
