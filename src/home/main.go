package home

import (
	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/types"
)

type HomeControllers struct {
	ServeHome types.Controller
}

func CreateHomeControllers() HomeControllers {
	return HomeControllers{
		ServeHome: auth.RedirectUnauthenticated(ServeHome),
	}
}
