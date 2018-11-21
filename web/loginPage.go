package web

import (
	"net/http"
	"strings"

	"jaredpearson.com/dbweb/data"
)

func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SessionStart(w, r)
	// TODO build login page for now just hardcode

	username := r.URL.Query()["username"]
	if len(username) != 1 {
		// TODO redirect to a login page instead
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := data.GetUserByUsername(strings.TrimSpace(username[0]))
	if err != nil {
		// TODO redirect to a login page instead
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	session.Set("username", user.Username())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
