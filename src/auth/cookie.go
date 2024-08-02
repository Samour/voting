package auth

import "net/http"

const cookieName = "SessionId"

func WriteSessionCookie(w http.ResponseWriter, sessionId string) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: sessionId,
		// Should set max-age, samesite etc for secure cookies
	})
}
