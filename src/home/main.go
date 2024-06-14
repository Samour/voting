package home

import "net/http"

type HomeControllers struct {
	ServeHome func(http.ResponseWriter, *http.Request)
}

func CreateHomeControllers() HomeControllers {
	return HomeControllers{
		ServeHome: ServeHome,
	}
}
