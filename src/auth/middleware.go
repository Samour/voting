package auth

import (
	"net/http"

	"github.com/Samour/voting/render"
	"github.com/Samour/voting/types"
)

const redirectAuthenticatedTarget = "/"
const redirectUnauthenticatedTarget = "/login"

func RedirectAuthenticated(c types.Controller) types.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r)
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(session.User.UserId) > 0 {
			http.Redirect(w, r, redirectAuthenticatedTarget, http.StatusFound)
			return
		}

		c(w, r)
	}
}

func RedirectUnauthenticated(c types.Controller) types.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r)
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(session.User.UserId) == 0 {
			http.Redirect(w, r, redirectUnauthenticatedTarget, http.StatusFound)
			return
		}

		c(w, r)
	}
}

func PreventUnauthenticated(c types.Controller) types.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r)
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(session.User.UserId) == 0 {
			render.ErrorPage(w, "Access Denied", http.StatusForbidden)
			return
		}

		c(w, r)
	}
}
