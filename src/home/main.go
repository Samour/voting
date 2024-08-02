package home

import (
	"github.com/Samour/voting/middleware"
	"github.com/Samour/voting/types"
)

type HomeControllers struct {
	ServeHome types.Controller
}

func CreateHomeControllers() HomeControllers {
	return HomeControllers{
		ServeHome: middleware.AuthenticatedWithRedirect(ServeHome),
	}
}
