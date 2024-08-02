package auth

import (
	"errors"
	"net/http"
	"time"
)

const cookieName = "SessionId"

func WriteSessionCookie(w http.ResponseWriter, sessionId string) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: sessionId,
		// Should set max-age, samesite etc for secure cookies
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Expires: time.Now().Add(-24 * time.Hour),
	})
}

func GetSessionId(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}

	return cookie.Value, nil
}

func GetSession(r *http.Request) (Session, error) {
	sessionId, err := GetSessionId(r)
	if err != nil || len(sessionId) == 0 {
		return Session{}, err
	}

	session, _ := RetrieveSession(sessionId)
	return session, nil
}
