package controllers

import (
	"html/template"
	"net/http"
)

const hot_reload = true

var tmplFunctions = template.FuncMap{
	"plus": func(i int, j int) int {
		return i + j
	},
}

type Renderer struct {
	HotReload bool
	globs     []string
	tmpl      *template.Template
}

func CreateRenderer(globs ...string) (*Renderer, error) {
	allGlobs := []string{"../resources/components/*.html"}
	allGlobs = append(allGlobs, globs...)

	tmpl, err := parseGlobs(allGlobs)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		HotReload: hot_reload,
		globs:     allGlobs,
		tmpl:      tmpl,
	}, nil
}

func parseGlobs(globs []string) (*template.Template, error) {
	template := template.New("index.html").Funcs(tmplFunctions)

	for _, glob := range globs {
		var err error
		template, err = template.ParseGlob(glob)
		if err != nil {
			return nil, err
		}
	}

	return template, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, model any) error {
	if r.HotReload {
		tmpl, err := parseGlobs(r.globs)
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, name, model)
	} else {
		return r.tmpl.ExecuteTemplate(w, name, model)
	}
}

func Must(r *Renderer, err error) *Renderer {
	if err != nil {
		panic(err.Error())
	}

	return r
}
