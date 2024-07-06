package render

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

type HttpResponse struct {
	HttpCode     int
	ErrorMessage string
	Model        any
}

type Renderer struct {
	HotReload bool
	globs     []string
	tmpl      *template.Template
}

type RenderAction struct {
	w        http.ResponseWriter
	r        *Renderer
	template string
}

func CreateRenderer(globs ...string) (*Renderer, error) {
	allGlobs := []string{"../resources/tmpl/common/*.html"}
	for _, glob := range globs {
		allGlobs = append(allGlobs, "../resources/tmpl/"+glob)
	}

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

func (r *Renderer) UsingTemplate(w http.ResponseWriter, name string) RenderAction {
	return RenderAction{
		w:        w,
		r:        r,
		template: name,
	}
}

func (a RenderAction) Render(response HttpResponse, err error) {
	if err != nil {
		ErrorPage(a.w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(response.ErrorMessage) > 0 {
		code := http.StatusBadRequest
		if response.HttpCode != 0 {
			code = response.HttpCode
		}
		ErrorPage(a.w, response.ErrorMessage, code)
		return
	}

	if response.HttpCode != 0 {
		a.w.WriteHeader(response.HttpCode)
	}
	err = a.r.Render(a.w, a.template, response.Model)
	if err != nil {
		ErrorPage(a.w, err.Error(), 0)
	}
}

func Must(r *Renderer, err error) *Renderer {
	if err != nil {
		panic(err.Error())
	}

	return r
}
