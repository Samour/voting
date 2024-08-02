package user

import (
	"net/http"

	"github.com/Samour/voting/user/login"
	"github.com/Samour/voting/user/signup"
)

type UserControllers struct {
	ServeSignUp  func(http.ResponseWriter, *http.Request)
	HandleSignUp func(http.ResponseWriter, *http.Request)

	ServeLogIn  func(http.ResponseWriter, *http.Request)
	HandleLogIn func(http.ResponseWriter, *http.Request)
}

func CreateUserControllers() UserControllers {
	return UserControllers{
		ServeSignUp:  signup.ServeSignUp,
		HandleSignUp: signup.HandleSignUp,

		ServeLogIn:  login.ServeLogIn,
		HandleLogIn: login.HandleLogIn,
	}
}
