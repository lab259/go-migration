[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/go-gormigrate/gormigrate/blob/master/LICENSE)
[![CircleCI](https://circleci.com/gh/lab259/go-migration/tree/master.svg?style=shield)](https://circleci.com/gh/lab259/go-migration/tree/master)
[![codecov](https://codecov.io/gh/lab259/go-migration/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/go-migration)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/go-migration)](https://goreportcard.com/report/github.com/lab259/go-migration)

# Migration

A simple migration framework for Golang.

## CLI

This migration does not have a binary CLI exported, but it has a built in
tool to be triggered as a CLI.

This approach was adopted because there is a implementation based on source
code implementation. For that to work we need to compile the tool.

## Usage

```bash
go get github.com/lab259/go-migration
```

Now we must implement the CLI:

```go
package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/lab259/go-migration"
	"github.com/lab259/go-migration/examples/mongo/db"
	_ "github.com/lab259/go-migration/examples/mongo/migrations" // Import all migrations
	"os"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost/test") // Starts the connection
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer session.Close() // Ensures closing the connection

	db.MongoSession = session

	source := migration.DefaultCodeSource()      // Default Source implementation is the CodeSource...
	reporter := migration.NewDefaultReporter()   // Default reporter implementation

	manager := migration.NewDefaultManager(migration.NewMongoDB(session), source)
	runner := migration.NewArgsRunner(reporter, manager, os.Exit) // Create a runner based on the arguments passed to the program
	runner.Run() // Run the command
}
```

Below is an example of a migration. The `NewCodeMigration` takes the file name
in consideration to extract the ID and Description of the migration. In the
following example the file name is "20171219012821_create_table_indexes.go".
So the ID will be 20171219012821 (but in time.Time) and the description will
be "create_table_indexes".

```go
package migrations

import (
	"github.com/globalsign/mgo"
	. "github.com/lab259/go-migration"
	"github.com/lab259/go-migration/examples/mongo/db"
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
```

Another thing the `NewCodeMigration` does is to auto register in the
`DefaultCodeSource`. If you create a migration through the `NewMigration`
you will need to register it manually.

## Motivation

At first, I was not intending to create my own migration framework until I got
stuck into two main problems.

1. I was using MongoDB and I could find a reliable framework to adopt;

2. Since most of the frameworks were based on SQL databases, SQL files or go files
   were used attached with `database/sql` package.

So, I decided I would create a minimalist framework that only cares about
retrieving and storing the current version and let the developer deal with any
database configurations and connections specifics.

As simple as it can be, it can track anything that can be versioned. From
NoSQLs, until text files using diffs.

## Mongo is NoSQL... Why bother migrations!?

Well... in the real world sometimes things change. And keep version of
the "API" in the record might not be the best approach in some cases.

However, the main reason for this implementation is to ensure indexes
are created correctly.

## Why it does't have a CLI command?

It could be implemented when you are using SQL files as source
(check the [source_directory.go](source_directory.go) implementation). However,
for the [source_code.go](source_code.go) you will need to create your own CLI.
Hence, all connections which must be open, and all migrations "classes" can be
compiled and linked together.

You can use the [`CLI`](cli.go) base implementation to extend and create your
own solution.

## Supported databases

- MongoDB (via [mgo](https://github.com/go-mgo/mgo))
- ~~MySQL~~ (TODO)
- PostgresSQL (via [lib/pq](https://github.com/lib/pq))
- ~~SQLite~~ (TODO)

# Other migration frameworks

- [Mattes Migrate](https://github.com/mattes/migrate)
- [Goose](https://github.com/pressly/goose)
- [Gemnasium Migrate](https://github.com/gemnasium/migrate)
- [GORM Migrate](https://github.com/go-gormigrate/gormigrate) (for the GORM)

## Extending

In order to add new technologies to work with this framework, you must implement
a `struct` implementing the [`Target`](target.go) interface.

The [`Target`](target.go) is a very simple interface responsible for getting and
setting the version for the database. You can check out how the
[`MongoDBTarget`](target_mongodb.go) was implemented at
[target_mongodb.go](target_mongodb.go).

## Bugs and features

For BUGs and new features feel free to create an [issue](issues).
