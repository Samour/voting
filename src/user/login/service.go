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

func logIn(username string, password string) (logInResult, render.HttpResponse, error) {
	credential, err := repository.LoadUsernamePasswordCredential(strings.ToLower(username))
	if err != nil {
		return logInResult{}, render.HttpResponse{}, err
	}

	// Always compare hash to prevent timing attacks
	passwordValid := verifyPassword(password, credential.PasswordHash)
	if len(credential.UserId) == 0 || !passwordValid {
		return logInResult{}, render.HttpResponse{
			HttpCode: http.StatusBadRequest,
			Model: LogInModel{
				ErrorMessage: "Username or password is incorrect",
				Username:     username,
			},
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

	return logInResult{
		SessionId: session.SessionId,
		Redirect:  "/",
	}, render.HttpResponse{}, nil
}

func verifyPassword(password string, passwordHash string) bool {
	var hashBytes = []byte(passwordHash)
	if len(passwordHash) == 0 {
		hashBytes = []byte("DummyPassword")
	}
	return bcrypt.CompareHashAndPassword(hashBytes, []byte(password)) == nil && len(passwordHash) > 0
}
