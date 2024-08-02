package signup

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/render"
	"github.com/Samour/voting/user/model"
	"github.com/Samour/voting/user/repository"
	"github.com/Samour/voting/utils"
	"golang.org/x/crypto/bcrypt"
)

type signUpResult struct {
	SessionId string
	Redirect  string
}

func createAccount(username string, password string) (signUpResult, render.HttpResponse, error) {
	validation := validateUsernamePassword(username, password)
	if len(validation) > 0 {
		return signUpResult{}, render.HttpResponse{
			HttpCode: http.StatusBadRequest,
			Model: signUpModel{
				Username:         username,
				ValidationErrors: validation,
			},
		}, nil
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return signUpResult{}, render.HttpResponse{}, err
	}
	user := model.User{
		UserId:      utils.IdGen(),
		DisplayName: username,
	}
	credentials := model.UsernamePasswordCredential{
		Username:     strings.ToLower(username),
		PasswordHash: string(passwordHashBytes),
		UserId:       user.UserId,
	}

	err = repository.InsertNewUser(user, credentials)
	if err != nil {
		if errors.Is(err, repository.UsernameUnavailableError{}) {
			return signUpResult{}, render.HttpResponse{
				HttpCode: http.StatusBadRequest,
				Model: signUpModel{
					ErrorMessage: "This username is not available",
					Username:     username,
				},
			}, nil
		}
		return signUpResult{}, render.HttpResponse{}, err
	}

	session := auth.CreateUserSession(auth.SessionUserDetails{
		UserId:      user.UserId,
		DisplayName: user.DisplayName,
		Username:    credentials.Username,
	})

	return signUpResult{
		SessionId: session.SessionId,
		Redirect:  "/",
	}, render.HttpResponse{}, nil
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
