package db

import "gopkg.in/mgo.v2"

var MongoSession *mgo.Session

// Mongo wraps a worker function to ensure that the session disposal after used.
func GetSession() *mgo.Session {
	return MongoSession.Clone()
}
