package types

import (
	"net/http"

	"github.com/Samour/voting/auth"
)

type Controller func(http.ResponseWriter, *http.Request)

type AuthenticatedController func(http.ResponseWriter, *http.Request, auth.Session)

type Middleware func(Controller) Controller
