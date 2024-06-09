package controllers

import "net/http"

func errorPage(w http.ResponseWriter, errorMsg string, httpCode int) {
	w.WriteHeader(httpCode)
	err := templates.ExecuteTemplate(w, "error.html", errorMsg)
	if err != nil {
		http.Error(w, errorMsg, httpCode)
	}
}
