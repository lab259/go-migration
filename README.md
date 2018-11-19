[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/go-gormigrate/gormigrate/blob/master/LICENSE)
[![CircleCI](https://circleci.com/gh/lab259/go-migration/tree/master.svg?style=shield)](https://circleci.com/gh/lab259/go-migration/tree/master)
[![codecov](https://codecov.io/gh/lab259/go-migration/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/go-migration)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/go-migration)](https://goreportcard.com/report/github.com/lab259/go-migration)

# Migration

A simple migration framework for Golang.

## Installation

```bash
go get github.com/lab259/go-migration
```

## Getting started

### The CLI

TODO

### Creating a migration

TODO

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

## Why it does't have a CLI command?

It could be implemented when you are using SQL files as source
(check the [source_directory.go](source_directory.go) implementation). However,
for the [source_code.go](source_code.go) you will need to create your own CLI.
Hence, all connections which must be open, and all migrations "classes" can be
compiled and linked together.

You can use the [`CLI`](cli.go) base implementation to extend and create your
own solution.

## Supported databases

* MongoDB (via [mgo](https://github.com/go-mgo/mgo))
* ~~MySQL~~ (TODO)
* ~~PostgresSQL~~ (TODO)
* ~~SQLite~~ (TODO)

# Other migration frameworks

* [Mattes Migrate](https://github.com/mattes/migrate)
* [Goose](https://github.com/pressly/goose)
* [Gemnasium Migrate](https://github.com/gemnasium/migrate)
* [GORM Migrate](https://github.com/go-gormigrate/gormigrate) (for the GORM)

## Extending

In order to add new technologies to work with this framework, you must implement
a `struct` implementing the [`Target`](target.go) interface.

The [`Target`](target.go) is a very simple interface responsible for getting and
setting the version for the database. You can check out how the
[`MongoDBTarget`](target_mongodb.go) was implemented at
[target_mongodb.go](target_mongodb.go).

## Bugs and features

For BUGs and new features feel free to create an [issue](issues).