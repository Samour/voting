package user

import (
	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/types"
	"github.com/Samour/voting/user/login"
	"github.com/Samour/voting/user/signup"
)

type UserControllers struct {
	ServeSignUp  types.Controller
	HandleSignUp types.Controller

	ServeLogIn  types.Controller
	HandleLogIn types.Controller
}

func CreateUserControllers() UserControllers {
	return UserControllers{
		ServeSignUp:  auth.RedirectAuthenticated(signup.ServeSignUp),
		HandleSignUp: signup.HandleSignUp,

		ServeLogIn:  auth.RedirectAuthenticated(login.ServeLogIn),
		HandleLogIn: login.HandleLogIn,
	}
}
