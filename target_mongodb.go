package migration

import "time"
import (
	"gopkg.in/mgo.v2"
)

// MongoDBTarget implements the migration.Target of the MongoDB.
//
// In order to get access to the MongoDB, migration.MongoDBTarget uses the MGo
// library (http://gopkg.in/mgo.v2).
type MongoDBTarget struct {
	session        *mgo.Session
	collectionName string
}

// mongoDBMigrationVersion represents the version stored on the MongoDB.
type mongoDBMigrationVersion struct {
	ID time.Time `bson:"_id"`
}

// NewMongoDB returns a new instance of the migration.MongoDBTarget
func NewMongoDB(session *mgo.Session) *MongoDBTarget {
	return &MongoDBTarget{
		collectionName: DefaultMigrationTable,
		session:        session,
	}
}

// Version implements the migration.Target.Version by fetching the current
// version of the database from the collection defined by
// migration.MongoDBTarget.SetCollectionName.
//
// It returns the current version of the database.
//
// Any error returned by the MGo, will be passed up to the caller.
func (t *MongoDBTarget) Version() (time.Time, error) {
	sess := t.session.Clone()
	defer sess.Close()

	c := sess.DB("").C(t.collectionName)
	var version mongoDBMigrationVersion
	q := c.Find(nil)
	if err := q.One(&version); err == nil {
		return version.ID.UTC(), nil
	} else if err == mgo.ErrNotFound {
		return NoVersion, nil
	} else {
		return NoVersion, err
	}
}

// SetVersion implements the migration.Target.SetVersion by storing the passed
// version on the database.
//
// It returns eny error returned by the MGo.
func (t *MongoDBTarget) SetVersion(id time.Time) error {
	sess := t.session.Clone()
	defer sess.Close()

	sess.DB("").C(t.collectionName).DropCollection()
	if id != NoVersion {
		if err := sess.DB("").C(t.collectionName).Insert(&mongoDBMigrationVersion{
			ID: id,
		}); err != nil {
			return err
		}
	}
	return nil
}

// SetCollectionName sets the name of the collection used to store the current
// version of the database.
func (t *MongoDBTarget) SetCollectionName(collection string) *MongoDBTarget {
	t.collectionName = collection
	return t
}

// Session returns the mgo.Session reference of this target.
func (t *MongoDBTarget) Session() *mgo.Session {
	return t.session
}
