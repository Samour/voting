package polls

import "net/http"

func ServeNewPoll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/polls/1/edit", http.StatusFound)
}
