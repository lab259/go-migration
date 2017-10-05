package migration

import "time"
import (
	"gopkg.in/mgo.v2"
	"reflect"
)

type MongoDBTarget struct {
	session        *mgo.Session
	collectionName string
}

type MongoMigrationMixin struct {
}

type MongoDBMigrationVersion struct {
	ID time.Time `bson:"_id"`
}

func NewMongoDB(session *mgo.Session) *MongoDBTarget {
	return &MongoDBTarget{
		collectionName: DEFAULT_MIGRATION_TABLE,
		session:        session,
	}
}

func (this *MongoMigrationMixin) Session() *mgo.Session {
	i := reflect.ValueOf(this).Interface()
	if m, ok := i.(Migration); ok {
		t := reflect.ValueOf(m.GetManager().Target()).Interface()
		if mongo, ok := t.(MongoDBTarget); ok {
			return mongo.Session()
		}
	}
	return nil
}

func (this *MongoDBTarget) Version() (time.Time, error) {
	sess := this.session.Clone()
	defer sess.Close()

	c := sess.DB("").C(this.collectionName)
	var version MongoDBMigrationVersion
	q := c.Find(nil)
	if err := q.One(&version); err == nil {
		return version.ID.UTC(), nil
	} else if err == mgo.ErrNotFound {
		return NOVERSION, nil
	} else {
		return NOVERSION, err
	}
}

func (this *MongoDBTarget) SetVersion(id time.Time) error {
	sess := this.session.Clone()
	defer sess.Close()

	sess.DB("").C(this.collectionName).DropCollection()
	if id != NOVERSION {
		if err := sess.DB("").C(this.collectionName).Insert(&MongoDBMigrationVersion{
			ID: id,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (this *MongoDBTarget) CollectionName(collection string) *MongoDBTarget {
	this.collectionName = collection
	return this
}

func (this *MongoDBTarget) Session() *mgo.Session {
	return this.session
}
