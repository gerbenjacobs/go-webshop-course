package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const cookieName = "flashes"

var store = sessions.NewCookieStore([]byte("this-really-should-be-a-secure-token-thats-not-stored-in-code"))

func storeAndSaveFlash(r *http.Request, w http.ResponseWriter, msg string) error {
	session, _ := store.Get(r, cookieName)
	session.AddFlash(msg)
	return session.Save(r, w)
}

func getFlashes(r *http.Request, w http.ResponseWriter) (map[string]string, error) {
	session, _ := store.Get(r, cookieName)
	flashes := session.Flashes()

	m := map[string]string{}
	for f := range flashes {
		fs := strings.SplitN(flashes[f].(string), "|", 2)
		if len(fs) == 2 {
			m[fs[0]] = fs[1]
		}
	}
	return m, session.Save(r, w)
}
