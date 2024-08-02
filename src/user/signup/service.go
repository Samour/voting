package signup

import (
	"errors"
	"net/http"

	"github.com/Samour/voting/render"
	"github.com/Samour/voting/user/model"
	"github.com/Samour/voting/user/repository"
	"github.com/Samour/voting/utils"
)

func createAccount(username string, password string) (*string, render.HttpResponse, error) {
	validation := validateUsernamePassword(username, password)
	if len(validation) > 0 {
		return nil, render.HttpResponse{
			HttpCode: http.StatusBadRequest,
			Model: SignUpModel{
				Username:         username,
				ValidationErrors: validation,
			},
		}, nil
	}

	user := model.User{
		UserId:      utils.IdGen(),
		DisplayName: username,
	}
	credentials := model.UsernamePasswordCredential{
		Username:     username,
		PasswordHash: password,
		UserId:       user.UserId,
	}

	err := repository.InsertNewUser(user, credentials)
	if err != nil {
		if errors.Is(err, repository.UsernameUnavailableError{}) {
			return nil, render.HttpResponse{
				HttpCode: http.StatusBadRequest,
				Model: SignUpModel{
					ErrorMessage: "This username is not available",
					Username:     username,
				},
			}, nil
		}
		return nil, render.HttpResponse{}, err
	}

	redirect := "/"
	return &redirect, render.HttpResponse{}, nil
}

func validateUsernamePassword(username string, password string) []string {
	hasLowerCase := false
	hasUpperCase := false
	hasSymbol := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' {
			hasLowerCase = true
		} else if c >= 'A' && c <= 'Z' {
			hasUpperCase = true
		} else {
			hasSymbol = true
		}
	}

	var validationErrors []string
	if len(username) < 3 {
		validationErrors = append(validationErrors, "Username must be at least 3 characters long")
	}
	if len(password) < 10 {
		validationErrors = append(validationErrors, "Password must be at least 10 characters long")
	}
	if !hasLowerCase {
		validationErrors = append(validationErrors, "Password must contain a lower case letter")
	}
	if !hasUpperCase {
		validationErrors = append(validationErrors, "Password must contain an upper case letter")
	}
	if !hasSymbol {
		validationErrors = append(validationErrors, "Password must containt a non-alphabetic symbol")
	}

	return validationErrors
}
