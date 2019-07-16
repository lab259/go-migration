package main

import (
	"os"

	"github.com/globalsign/mgo"
	rlog2 "github.com/lab259/rlog/v2"

	"github.com/lab259/go-migration"
	"github.com/lab259/go-migration/examples/mongo/db"
	_ "github.com/lab259/go-migration/examples/mongo_rlog/migrations"
	"github.com/lab259/go-migration/rlog"
)

func main() {
	logger := rlog2.WithFields(nil)

	session, err := mgo.Dial("mongodb://localhost/test")
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	defer session.Close()

	db.MongoSession = session
	source := migration.DefaultCodeSource()
	reporter := rlog.NewRLogReporter(logger, os.Exit)

	manager := migration.NewDefaultManager(migration.NewMongoDB(session.DB("")), source)
	runner := migration.NewArgsRunner(reporter, manager, os.Exit)
	runner.Run(nil)
}
