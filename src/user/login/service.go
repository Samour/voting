package login

import (
	"net/http"
	"strings"

	"github.com/Samour/voting/render"
	"github.com/Samour/voting/user/repository"
	"golang.org/x/crypto/bcrypt"
)

func logIn(username string, password string) (*string, render.HttpResponse, error) {
	credential, err := repository.LoadUsernamePasswordCredential(strings.ToLower(username))
	if err != nil {
		return nil, render.HttpResponse{}, err
	}

	// Always compare hash to prevent timing attacks
	passwordValid := verifyPassword(password, credential.PasswordHash)
	if len(credential.UserId) == 0 || !passwordValid {
		return nil, render.HttpResponse{
			HttpCode: http.StatusBadRequest,
			Model: LogInModel{
				ErrorMessage: "Username or password is incorrect",
				Username:     username,
			},
		}, nil
	}

	redirect := "/"
	return &redirect, render.HttpResponse{}, nil
}

func verifyPassword(password string, passwordHash string) bool {
	var hashBytes = []byte(passwordHash)
	if len(passwordHash) == 0 {
		hashBytes = []byte("DummyPassword")
	}
	return bcrypt.CompareHashAndPassword(hashBytes, []byte(password)) == nil && len(passwordHash) > 0
}
