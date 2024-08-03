package login

import (
	"net/http"
	"strings"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/render"
	"github.com/Samour/voting/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type logInResult struct {
	SessionId string
	Redirect  string
}

func logIn(username string, password string, redirect string) (logInResult, render.HttpResponse, error) {
	if len(username) == 0 {
		return logInResult{}, render.HttpResponse{
			HttpCode: http.StatusBadRequest,
			Model:    buildLogInModel(username, redirect, "Must provide a username"),
		}, nil
	}

	credential, err := repository.LoadUsernamePasswordCredential(strings.ToLower(username))
	if err != nil {
		return logInResult{}, render.HttpResponse{}, err
	}

	// Always compare hash to prevent timing attacks
	passwordValid := verifyPassword(password, credential.PasswordHash)
	if len(credential.UserId) == 0 || !passwordValid {
		return logInResult{}, render.HttpResponse{
			HttpCode: http.StatusBadRequest,
			Model:    buildLogInModel(username, redirect, "Username or password is incorrect"),
		}, nil
	}

	user, err := repository.LoadUser(credential.UserId)
	if err != nil {
		return logInResult{}, render.HttpResponse{}, err
	}

	session := auth.CreateUserSession(auth.SessionUserDetails{
		UserId:      credential.UserId,
		DisplayName: user.DisplayName,
		Username:    credential.Username,
	})

	if len(redirect) == 0 {
		redirect = "/"
	}
	return logInResult{
		SessionId: session.SessionId,
		Redirect:  redirect,
	}, render.HttpResponse{}, nil
}

func verifyPassword(password string, passwordHash string) bool {
	var hashBytes = []byte(passwordHash)
	if len(passwordHash) == 0 {
		hashBytes = []byte("DummyPassword")
	}
	return bcrypt.CompareHashAndPassword(hashBytes, []byte(password)) == nil && len(passwordHash) > 0
}
