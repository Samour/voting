package controllers

import "net/http"

func ErrorPage(w http.ResponseWriter, errorMsg string, httpCode int) {
	w.WriteHeader(httpCode)
	err := Templates.ExecuteTemplate(w, "error.html", errorMsg)
	if err != nil {
		http.Error(w, errorMsg, httpCode)
	}
}
