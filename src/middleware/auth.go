package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/render"
	"github.com/Samour/voting/types"
)

const redirectAuthenticatedTarget = "/"
const redirectUnauthenticatedTarget = "/login"
const redirectParam = "redirect"

func Unauthenticated(c types.Controller) types.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSession(r)
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

func AuthenticatedWithRedirect(c types.AuthenticatedController) types.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSession(r)
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(session.User.UserId) == 0 {
			redirect := fmt.Sprintf("%s?%s=%s", redirectUnauthenticatedTarget, redirectParam, url.QueryEscape(r.URL.Path))
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		c(w, r, session)
	}
}

func AuthenticatedWithError(c types.AuthenticatedController) types.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSession(r)
		if err != nil {
			render.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(session.User.UserId) == 0 {
			render.ErrorPage(w, "Access Denied", http.StatusForbidden)
			return
		}

		c(w, r, session)
	}
}

func GetAuthRedirect(r *http.Request) string {
	return r.Form.Get(redirectParam)
}
