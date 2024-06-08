package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var count = 0

func main() {
	http.HandleFunc("/increment", handleIncrement)
	http.HandleFunc("/", handleRoot)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		serveErrorPage(w, "Not found", http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("resources/pages/home.html")
	if err != nil {
		serveErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, count)
	if err != nil {
		serveErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIncrement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		serveErrorPage(w, "Not found", http.StatusNotFound)
		return
	}

	count++
	fmt.Fprintf(w, "%d", count)
}

func serveErrorPage(w http.ResponseWriter, errorMsg string, httpCode int) {
	http.Error(w, errorMsg, httpCode)
}
