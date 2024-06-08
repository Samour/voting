package home

import (
	"net/http"

	"github.com/Samour/voting/controllers"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	err := controllers.Templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		controllers.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
