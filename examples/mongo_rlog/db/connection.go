package db

import "github.com/globalsign/mgo"

// MongoSession is the global Session of the mongo connection
var MongoSession *mgo.Session

// GetSession returns a clone of the current global session
func GetSession() *mgo.Session {
	return MongoSession.Clone()
}
