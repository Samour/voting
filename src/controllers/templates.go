package controllers

import (
	"html/template"
	"net/http"
)

const hot_reload = true

var templates = template.Must(parseFiles())

var tmplFunctions = template.FuncMap{
	"plus": func(i int, j int) int {
		return i + j
	},
}

func parseFiles() (*template.Template, error) {
	return template.New("home.html").Funcs(tmplFunctions).ParseFiles(
		"../resources/components/edit_poll_options.html",
		"../resources/components/page_footer.html",
		"../resources/components/page_header.html",
		"../resources/components/view_poll_navigation.html",

		"../resources/pages/edit_poll.html",
		"../resources/pages/error.html",
		"../resources/pages/home.html",
		"../resources/pages/view_poll.html",
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
