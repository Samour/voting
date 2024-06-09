package controllers

import (
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		errorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
