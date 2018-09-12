package migrations

import (
	. "../../.."
	"../db"
	"github.com/globalsign/mgo"
)

func init() {
	NewCodeMigration(
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
		},
	)
}
