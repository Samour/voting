package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/home"
	"github.com/Samour/voting/polls"
	"github.com/Samour/voting/user"
)

func main() {
	clients.WarmDynamoDbClient()

	homeControllers := home.CreateHomeControllers()
	userControllers := user.CreateUserControllers()
	pollControllers := polls.CreatePollControllers()

	static := http.FileServer(http.Dir("../resources/static/"))
	http.Handle("GET /static/", http.StripPrefix("/static/", static))

	http.HandleFunc("GET /{$}", homeControllers.ServeHome)

	http.HandleFunc("GET /signup", userControllers.ServeSignUp)
	http.HandleFunc("POST /signup", userControllers.HandleSignUp)

	http.HandleFunc("GET /login", userControllers.ServeLogIn)
	http.HandleFunc("POST /login", userControllers.HandleLogIn)
	http.HandleFunc("GET /logout", userControllers.ServeLogOut)

	http.HandleFunc("GET /polls/{id}/{$}", pollControllers.ServeViewPoll)
	http.HandleFunc("PUT /polls/{id}/status/{$}", pollControllers.HandlePollStatusChange)

	http.HandleFunc("GET /polls/new/{$}", pollControllers.ServeNewPoll)

	http.HandleFunc("GET /polls/{id}/edit/{$}", pollControllers.ServeEditPoll)
	http.HandleFunc("POST /polls/{id}/{$}", pollControllers.ServeSavePoll)
	http.HandleFunc("PATCH /polls/options/{$}", pollControllers.HandlePatchPoll)

	http.HandleFunc("GET /polls/{id}/vote/{$}", pollControllers.ServeVotePoll)
	http.HandleFunc("POST /polls/{id}/vote/fptp/{$}", pollControllers.HandleCastFptpVote)
	http.HandleFunc("PATCH /polls/{id}/vote/rankedchoice/{$}", pollControllers.HandlePatchRankedChoice)
	http.HandleFunc("POST /polls/{id}/vote/rankedchoice/{$}", pollControllers.HandleCastRankedChoiceVote)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
