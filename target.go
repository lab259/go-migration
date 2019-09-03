package migration

import "time"

// NoVersion represents a zero version
var NoVersion = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

// Target describes the common interface for target of migrations. Each target
// is a implementation of a specific database (or anything versionable).
type Target interface {
	// Version returns the current version of the database.
	Version() (time.Time, error)

	// AddMigration persists the version on the database (or any other mean
	// necessary).
	AddMigration(summary *Summary) error

	// RemoveMigration removes a migration record from the database.
	RemoveMigration(summary *Summary) error

	// MigrationsExecuted returns all the migrations executed.
	MigrationsExecuted() ([]time.Time, error)
}

// BeforeRun describes a hook to be called before the Runner actually run.
type BeforeRun interface {
	BeforeRun(executionContext interface{})
}
