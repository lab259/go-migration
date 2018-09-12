package main

import (
	"../.."
	"../../examples/mongo/db"
	_ "../../examples/mongo/migrations"
	"fmt"
	"github.com/globalsign/mgo"
	"os"
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
	runner := migration.NewArgsRunner(reporter, manager, os.Exit)
	runner.Run()
}
