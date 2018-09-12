package db

import "github.com/globalsign/mgo"

var MongoSession *mgo.Session

// Mongo wraps a worker function to ensure that the session disposal after used.
func GetSession() *mgo.Session {
	return MongoSession.Clone()
}
