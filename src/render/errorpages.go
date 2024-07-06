package render

import "net/http"

var errorRenderer = Must(createErrorRenderer())

func createErrorRenderer() (*Renderer, error) {
	renderer, err := CreateRenderer("pages/error.html")
	if err != nil {
		return nil, err
	}

	renderer.HotReload = false
	return renderer, nil
}

func ErrorPage(w http.ResponseWriter, errorMsg string, httpCode int) {
	if httpCode != 0 {
		w.WriteHeader(httpCode)
	}
	err := errorRenderer.Render(w, "error.html", errorMsg)
	if err != nil {
		http.Error(w, errorMsg, httpCode)
	}
}
