package migrations

import (
	"gopkg.in/mgo.v2"
	. "github.com/jamillosantos/migration"
	"github.com/jamillosantos/migration/examples/mongo/db"
)

func init() {
	Register(
		NewMigration(
			NewMigrationId("20171219012821"),
			"create customer indexes",
			func() error {
				// Get the connection reference
				session := db.GetSession()
				defer session.Close()
				//-------------------

				c := session.DB("").C("customers")
				err := c.EnsureIndex(mgo.Index{
					Name: "NameIndex",
					Key:  []string{"name"},
				})
				return err // Return the error of the operation
			}, func() error {
				session := db.GetSession()
				defer session.Close()

				err := session.DB("").C("customers").DropIndexName("NameIndex")
				return err
			}))
}
