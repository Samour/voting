package controllers

import "net/http"

func ErrorPage(w http.ResponseWriter, errorMsg string, httpCode int) {
	http.Error(w, errorMsg, httpCode)
}
