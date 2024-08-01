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
