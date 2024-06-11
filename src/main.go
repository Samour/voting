package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/controllers"
)

func main() {
	clients.WarmDynamoDbClient()

	static := http.FileServer(http.Dir("resources/static/"))
	http.Handle("GET /static/", http.StripPrefix("/static/", static))
	http.HandleFunc("GET /{$}", controllers.ServeHome)
	http.HandleFunc("GET /polls/{id}/{$}", controllers.ServeViewPoll)
	http.HandleFunc("GET /polls/new/{$}", controllers.ServeNewPoll)
	http.HandleFunc("GET /polls/{id}/edit/{$}", controllers.ServeEditPoll)
	http.HandleFunc("POST /polls/{id}/{$}", controllers.ServeSavePoll)
	http.HandleFunc("PATCH /polls/options/{$}", controllers.HandlePatchPoll)
	http.HandleFunc("PUT /polls/{id}/status/{$}", controllers.HandlePollStatusChange)
	http.HandleFunc("GET /polls/{id}/vote/{$}", controllers.ServeVotePoll)
	http.HandleFunc("POST /polls/{id}/vote/{$}", controllers.ServeCastVote)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
