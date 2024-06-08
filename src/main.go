package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Samour/voting/controllers/home"
	"github.com/Samour/voting/controllers/polls"
)

func main() {
	static := http.FileServer(http.Dir("resources/static/"))
	http.Handle("GET /static/", http.StripPrefix("/static/", static))
	http.HandleFunc("GET /{$}", home.ServeHome)
	http.HandleFunc("GET /polls/new", polls.ServeNewPoll)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
