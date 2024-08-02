package auth

import (
	"sync"

	"github.com/Samour/voting/utils"
)

type Session struct {
	SessionId string
	User      SessionUserDetails
}

type SessionUserDetails struct {
	UserId      string
	DisplayName string
	Username    string
}

var sessionStore = sync.Map{}

func CreateUserSession(userDetails SessionUserDetails) Session {
	sessionId := utils.IdGenOfLength(20)
	session := Session{
		SessionId: sessionId,
		User:      userDetails,
	}

	sessionStore.Store(sessionId, session)

	return session
}
