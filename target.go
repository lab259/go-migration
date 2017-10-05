package migration

import "time"

// NoVersion represents a zero version
var NoVersion = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

// Target describes the common interface for target of migrations. Each target
// is a implementation of a specific database (or anything versionable).
type Target interface {
	// Version returns the current version of the database.
	Version() (time.Time, error)

	// SetVersion persists the version on the database (or any other mean
	// necessary).
	SetVersion(version time.Time) error
}
