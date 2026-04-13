package handlers

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitSessionStore() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 12,
		HttpOnly: true,
	}
}

func getSessionUserID(r *http.Request) int {
	session, _ := store.Get(r, "session")
	id, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}
	return id
}

func getSessionRole(r *http.Request) string {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok {
		return ""
	}
	return role
}

func setFlash(w http.ResponseWriter, r *http.Request, message string) {
	session, _ := store.Get(r, "session")
	session.Values["flash"] = message
	session.Save(r, w)
}

func getFlash(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, "session")
	msg, ok := session.Values["flash"].(string)
	if !ok {
		return ""
	}
	delete(session.Values, "flash")
	session.Save(r, w)
	return msg
}
