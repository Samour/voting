package types

import "net/http"

type Controller func(http.ResponseWriter, *http.Request)

type Middleware func(Controller) Controller
