package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"jaredpearson.com/dbweb/data"

	"jaredpearson.com/dbweb/web"
)

func determinePort() string {
	port := os.Getenv("PORT")
	if len(port) != 0 {
		return port
	}

	// default port
	return "8080"
}

type HomePageData struct {
	pageTitle string
	userInfo  web.UserInfo
	Sets      []*data.MiniatureSet
}

func (page HomePageData) PageTitle() string {
	return page.pageTitle
}
func (page HomePageData) UserInfo() web.UserInfo {
	return page.userInfo
}

func showHome(w http.ResponseWriter, r *http.Request) {
	userInfo, _ := web.UserInfoFromRequest(r)
	sets, err := data.GetMiniatureSets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := HomePageData{
		pageTitle: "",
		userInfo:  userInfo,
		Sets:      sets,
	}
	web.ShowTemplateInMainLayout(w, r, "home", data)
}

// fillRequestSession puts the Session in the request context under the
// token sessionRequestToken
func fillRequestSession(sessionManager *web.SessionManager) web.HttpMiddleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sessionManager.HasSession(r) {
				session, err := sessionManager.ReadSession(r)
				if err == nil {
					newContext := context.WithValue(r.Context(), web.SessionRequestToken, session)
					r = r.WithContext(newContext)
				}
			}
			handler.ServeHTTP(w, r)
		})
	}
}

// fillUser will retrieve the user information from the sesssion
// and populate the user information in the Request context if
// the user is logged in. If the user is not logged in, then Request
// context will remain unchanged.
//
// This requires the session to be populated. See fillRequestSession
func fillUserMiddleware() web.HttpMiddleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := r.Context().Value(web.SessionRequestToken)
			if s != nil {
				username := s.(web.Session).Get("username").(string)
				if len(username) > 0 {
					// TODO attempt to load the user name
					user := web.UserInfo{
						Username: username,
					}
					newContext := context.WithValue(r.Context(), web.AuthUserToken, user)
					r = r.WithContext(newContext)
				}
			}
			handler.ServeHTTP(w, r)
		})
	}
}

var sessionManager *web.SessionManager

func init() {
	sessionProvider := web.NewMongoDbSessionProvider(
		data.GetMongoSession,
		data.GetMongoDbName(),
		"sessions",
	)
	sessionProvider.InitializeMongoDb()
	sessionManager, _ = web.NewSessionManager("dbsession", sessionProvider)
}

func main() {

	fillSession := fillRequestSession(sessionManager)
	fillUser := fillUserMiddleware()
	mwChain := web.ChainMiddleware(fillSession, fillUser)

	http.Handle("/", mwChain(http.HandlerFunc(showHome)))
	http.Handle("/login", mwChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.SessionStart(w, r)

		// TODO build login page for now just hardcode
		session.Set("username", "user@dreamblade.com")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})))
	http.Handle("/miniature/", mwChain(http.HandlerFunc(web.ShowMiniatureDetailPage)))
	http.Handle("/set/", mwChain(http.HandlerFunc(web.ShowSetDetailPage)))

	port := determinePort()
	log.Printf("Server started on %s", port)

	// TODO fix issue with missing pages being sent to home
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
