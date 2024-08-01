package user

import (
	"net/http"

	"github.com/Samour/voting/user/signup"
)

type UserControllers struct {
	ServeSignUp func(http.ResponseWriter, *http.Request)
}

func CreateUserControllers() UserControllers {
	return UserControllers{
		ServeSignUp: signup.ServeSignUp,
	}
}
