package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const CACHE_TEMPLATES = true

var count = 0

var templates = template.Must(template.ParseFiles(
	"resources/components/page_footer.html",
	"resources/components/page_header.html",

	"resources/pages/home.html",
	"resources/pages/new_poll.html",
))

func main() {
	static := http.FileServer(http.Dir("resources/static/"))
	http.Handle("GET /static/", http.StripPrefix("/static/", static))
	http.HandleFunc("GET /{$}", serveHome)
	http.HandleFunc("GET /polls/new", serveNewPoll)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "home.html", count)
	if err != nil {
		serveErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveNewPoll(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "new_poll.html", nil)
	if err != nil {
		serveErrorPage(w, err.Error(), http.StatusInternalServerError)
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

func serveErrorPage(w http.ResponseWriter, errorMsg string, httpCode int) {
	http.Error(w, errorMsg, httpCode)
}
