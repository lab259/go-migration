package migrations

import (
	"github.com/globalsign/mgo"

	. "github.com/lab259/go-migration"
	"github.com/lab259/go-migration/examples/mongo/db"
)

func init() {
	NewCodeMigration(
		func(executionContext interface{}) error {
			// Get the connection reference
			session := db.GetSession()
			defer session.Close()
			// -------------------

			c := session.DB("").C("customers")
			err := c.EnsureIndex(mgo.Index{
				Name: "NameIndex",
				Key:  []string{"name"},
			})
			return err // Return the error of the operation
		}, func(executionContext interface{}) error {
			session := db.GetSession()
			defer session.Close()

			err := session.DB("").C("customers").DropIndexName("NameIndex")
			return err
		},
	)
}
