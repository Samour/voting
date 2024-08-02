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

func RetrieveSession(sessionId string) (Session, bool) {
	item, exists := sessionStore.Load(sessionId)
	if !exists {
		return Session{}, false
	}

	if session, ok := item.(Session); ok {
		return session, true
	} else {
		return Session{}, false
	}
}

func RemoveSession(sessionId string) {
	sessionStore.Delete(sessionId)
}
