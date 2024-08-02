package user

import (
	"github.com/Samour/voting/middleware"
	"github.com/Samour/voting/types"
	"github.com/Samour/voting/user/login"
	"github.com/Samour/voting/user/signup"
)

type UserControllers struct {
	ServeSignUp  types.Controller
	HandleSignUp types.Controller

	ServeLogIn  types.Controller
	HandleLogIn types.Controller
	ServeLogOut types.Controller
}

func CreateUserControllers() UserControllers {
	return UserControllers{
		ServeSignUp:  middleware.RedirectAuthenticated(signup.ServeSignUp),
		HandleSignUp: signup.HandleSignUp,

		ServeLogIn:  middleware.RedirectAuthenticated(login.ServeLogIn),
		HandleLogIn: login.HandleLogIn,
		ServeLogOut: middleware.AuthenticatedWithRedirect(login.ServeLogOut),
	}
}
