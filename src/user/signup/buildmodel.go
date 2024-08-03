package signup

import (
	"fmt"
	"net/url"

	"github.com/Samour/voting/middleware"
)

func buildSignUpModel(username string, redirect string, errorMessage string, validationErrors []string) signUpModel {
	logInUrl := "/login"
	if len(redirect) > 0 {
		logInUrl = fmt.Sprintf("%s?%s=%s", logInUrl, middleware.ParamRedirect, url.QueryEscape(redirect))
	}

	return signUpModel{
		Username:         username,
		ErrorMessage:     errorMessage,
		ValidationErrors: validationErrors,
		LogInUrl:         logInUrl,
	}
}
