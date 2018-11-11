package web

// Web server session managment and storage. The implementation is very
// similar to the following post:
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/06.2.html
//
// This implementation is basic and *does not* protect against session hijacking
//
// TODO expiration of cookies

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Session interface {
	SessionID() string
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string) error
}

type SessionProvider interface {
	InitializeSession(sid string) (Session, error)
	ReadSession(sid string) (Session, error)
}

// MongoDbSession is a web session backed by MongoDb
type MongoDbSession struct {
	provider  *MongoDbSessionProvider
	sessionID string
	data      map[string]interface{}
}

func (session *MongoDbSession) SessionID() string {
	return session.sessionID
}
func (session *MongoDbSession) Set(key string, value interface{}) error {
	session.data[key] = value
	return session.provider.UpdateSession(session)
}
func (session *MongoDbSession) Get(key string) interface{} {
	return session.data[key]
}
func (session *MongoDbSession) Delete(key string) error {
	delete(session.data, key)
	return nil
}

type CreateMongoDbSession func() *mgo.Session

// MongoDbSessionProvider is a SessionProvider that's backed by MongoDb.
type MongoDbSessionProvider struct {
	lock                 sync.Mutex
	createMongoDbSession CreateMongoDbSession
	databaseName         string
	collectionName       string
}

func NewMongoDbSessionProvider(
	createMongoDbSession CreateMongoDbSession,
	databaseName,
	collectionName string) *MongoDbSessionProvider {
	return &MongoDbSessionProvider{
		createMongoDbSession: createMongoDbSession,
		databaseName:         databaseName,
		collectionName:       collectionName,
	}
}

func (provider *MongoDbSessionProvider) InitializeMongoDb() (err error) {
	session := provider.createMongoDbSession()
	defer session.Close()

	err = session.DB(provider.databaseName).C(provider.collectionName).EnsureIndex(mgo.Index{
		Key: []string{"sid"},
	})
	return
}

func (provider *MongoDbSessionProvider) InitializeSession(sid string) (Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	session := &MongoDbSession{
		sessionID: sid,
		provider:  provider,
		data:      make(map[string]interface{}),
	}
	return session, nil
}
func (provider *MongoDbSessionProvider) ReadSession(sid string) (Session, error) {
	session := provider.createMongoDbSession()
	defer session.Close()

	var sessionData map[string]interface{}

	collection := session.DB(provider.databaseName).C(provider.collectionName)

	query := bson.M{"sid": sid}

	err := collection.Find(query).One(&sessionData)
	if err != nil {
		return nil, err
	}

	// assume v1 of session data doc
	data := sessionData["data"].(map[string]interface{})

	return &MongoDbSession{
		sessionID: sid,
		provider:  provider,
		data:      data,
	}, nil
}
func (provider *MongoDbSessionProvider) UpdateSession(session *MongoDbSession) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	mongoSession := provider.createMongoDbSession()
	defer mongoSession.Close()

	doc := make(map[string]interface{})
	doc["sid"] = session.SessionID()
	doc["version"] = 1
	doc["data"] = session.data

	query := bson.M{"sid": session.SessionID()}

	collection := mongoSession.DB(provider.databaseName).C(provider.collectionName)
	_, err := collection.Upsert(query, doc)

	return err
}

// SessionManager is used by the application to manage sessions
type SessionManager struct {
	cookieName string
	lock       sync.Mutex
	provider   SessionProvider
}

// NewSessionManager creates a new session manager. This should only be
// invoked once when the application is started.
// cookiName is the name used to store the session information
// sessionProvider is where the sessions are to be stored
func NewSessionManager(cookieName string, provider SessionProvider) (*SessionManager, error) {
	return &SessionManager{
		cookieName: cookieName,
		provider:   provider,
	}, nil
}

func (manager *SessionManager) generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// createNewSession is a private function used to create a new session. The client
// should use SessionStart to create a session.
func (manager *SessionManager) createNewSession(w http.ResponseWriter, r *http.Request) (session Session) {
	sid := manager.generateSessionID()
	session, _ = manager.provider.InitializeSession(sid)
	cookie := http.Cookie{
		Name:     manager.cookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	}
	http.SetCookie(w, &cookie)
	return session
}

// SessionStart should be called by the client to initialize a new session or
// retrieve an existing session (from the ID in the cookie)
func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	var err error

	// attempt to read the session from the request
	session, err = manager.ReadSession(r)
	if err != nil || session == nil {
		// create a new session if we didn't find one
		session = manager.createNewSession(w, r)
	}
	return
}

func (manager *SessionManager) ReadSession(r *http.Request) (session Session, err error) {
	var cookie *http.Cookie
	cookie, err = r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		// pass the error back to the caller
		return
	}
	sid, _ := url.QueryUnescape(cookie.Value)
	session, err = manager.provider.ReadSession(sid)
	return
}

// HasSession determines if a session has already been associated to the request
func (manager *SessionManager) HasSession(r *http.Request) bool {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)

	// TODO make sure the session hasn't expired
	return err == nil && cookie.Value != ""
}
