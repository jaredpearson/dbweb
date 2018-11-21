package data

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

// ErrUserNotFound error returned when a user could not be found
type ErrUserNotFound struct {
	username string
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user not found: %s", e.username)
}

func NewErrUserNotFound(username string) *ErrUserNotFound {
	return &ErrUserNotFound{
		username,
	}
}

const (
	mongoUserCollectionName = "users"
	userDocCurrentVersion   = 1
)

type User interface {
	Username() string
}

type user struct {
	username string
}

func (u user) Username() string {
	return u.username
}

type userDto struct {
	Version  int
	Username string
}

func (userData userDto) toUser() User {
	return &user{
		username: userData.Username,
	}
}

func GetUserByUsername(username string) (user User, err error) {
	mongoSession := GetMongoSession()
	defer mongoSession.Close()

	query := bson.M{"username": username}
	userData := userDto{}

	userCollection := session.DB(GetMongoDbName()).C(mongoUserCollectionName)
	err = userCollection.Find(query).One(&userData)
	if err != nil {
		err = NewErrUserNotFound(username)
		return
	}
	user = userData.toUser()
	return
}

func AddUser(username string) (err error) {
	mongoSession := GetMongoSession()
	defer mongoSession.Close()

	user := userDto{
		userDocCurrentVersion,
		username,
	}

	// TODO make sure the username isn't already taken

	userCollection := session.DB(GetMongoDbName()).C(mongoUserCollectionName)
	err = userCollection.Insert(user)
	return
}
