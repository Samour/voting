package polls

import (
	"net/http"

	"github.com/Samour/voting/controllers"
)

func ServeNewPoll(w http.ResponseWriter, r *http.Request) {
	err := controllers.Templates.ExecuteTemplate(w, "new_poll.html", nil)
	if err != nil {
		controllers.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
