package main

import (
	"github.com/jamillosantos/migration"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/jamillosantos/migration/examples/mongo/db"
	_ "github.com/jamillosantos/migration/examples/mongo/migrations"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost/test")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer session.Close()

	db.MongoSession = session

	source := migration.DefaultCodeSource()
	reporter := migration.NewDefaultReporter()

	manager := migration.NewDefaultManager(migration.NewMongoDB(session), source)
	runner := migration.NewArgsRunner(reporter, manager)
	runner.Run()
}
