package login

import (
	"fmt"
	"net/url"

	"github.com/Samour/voting/middleware"
)

func buildLogInModel(username string, redirect string, errorMessage string) logInModel {
	signUpUrl := "/signup"
	if len(redirect) > 0 {
		signUpUrl = fmt.Sprintf("%s?%s=%s", signUpUrl, middleware.ParamRedirect, url.QueryEscape(redirect))
	}

	return logInModel{
		ErrorMessage: errorMessage,
		Username:     username,
		SignUpUrl:    signUpUrl,
	}
}
