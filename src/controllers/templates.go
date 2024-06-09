package controllers

import (
	"html/template"
	"net/http"
)

const hot_reload = true

var templates = template.Must(parseFiles())

func parseFiles() (*template.Template, error) {
	return template.ParseFiles(
		"../resources/components/page_footer.html",
		"../resources/components/page_header.html",

		"../resources/pages/edit_poll.html",
		"../resources/pages/error.html",
		"../resources/pages/home.html",
	)
}

func renderTemplate(w http.ResponseWriter, name string, model any) error {
	if hot_reload {
		tmpl, err := parseFiles()
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, name, model)
	} else {
		return templates.ExecuteTemplate(w, name, model)
	}
}
