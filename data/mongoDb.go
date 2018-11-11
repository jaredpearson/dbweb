package data

import (
	"log"

	"github.com/globalsign/mgo"
)

var session *mgo.Session

func GetMongoDbName() string {
	return "dreamblade"
}

func GetMongoSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.Dial("mongodb://127.0.0.1:27017")
		if err != nil {
			log.Fatalf("Unable to contact MongoDB\n\t%v", err)
		}

		cred := mgo.Credential{
			Username: "mongoadmin",
			Password: "secret",
		}
		err = session.Login(&cred)
		if err != nil {
			log.Fatalf("Unable to authenticate\n\t%v", err)
		}
	}
	return session.Clone()
}
