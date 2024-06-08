package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const CACHE_TEMPLATES = false

var count = 0

var templates = template.Must(template.ParseFiles("resources/pages/home.html"))

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

	err := renderTemplate(w, "home.html", count)
	if err != nil {
		serveErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderTemplate(w http.ResponseWriter, t string, data any) error {
	if CACHE_TEMPLATES {
		return templates.ExecuteTemplate(w, t, data)
	} else {
		tmpl, err := template.ParseFiles("resources/pages/" + t)
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data)
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
