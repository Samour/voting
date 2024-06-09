package controllers

import (
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
