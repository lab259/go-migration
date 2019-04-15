package migration

import (
	"fmt"
	"time"
)

// Direction is the enum that represents the direction of the migration.
//
// Its possible values are `MigrationDirectionDo` and `MigrationDirectionUndo`.
type Direction uint

const (
	// DirectionDo is the forward direction.
	DirectionDo Direction = iota
	// DirectionUndo is the backwards direction.
	DirectionUndo Direction = iota
)

const migrationIDFormat = "20060102150405"

// DefaultMigrationTable is the default name of the migrations table.
const DefaultMigrationTable = "_migrations"

// NewMigrationID creates a new ID from a string.
func NewMigrationID(str string) time.Time {
	result, err := time.Parse(migrationIDFormat, str)
	if err == nil {
		return result
	}
	panic(fmt.Sprintf("%s is not an valid ID. (%s)", str, err))
}

// Migration is the interface that describes the common behavior that a
// migration should have to be manageable by the migration.Manager.
type Migration interface {
	GetID() time.Time
	GetDescription() string
	Do(executionContext interface{}) error
	Undo(executionContext interface{}) error
	GetManager() Manager
	SetManager(manager Manager) Migration
}
