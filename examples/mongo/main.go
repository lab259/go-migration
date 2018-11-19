package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/lab259/go-migration"
	"github.com/lab259/go-migration/examples/mongo/db"
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
