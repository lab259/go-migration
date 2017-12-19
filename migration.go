package migration

import (
	"time"
	"fmt"
)

type MigrationDirection uint

const (
	MigrationDirectionDo   MigrationDirection = iota
	MigrationDirectionUndo MigrationDirection = iota
)

const MigrationIdFormat = "20060102150405"

// DefaultMigrationTable is the default name of the migrations table.
const DefaultMigrationTable = "_migrations"

func NewMigrationId(str string) time.Time {
	result, err := time.Parse(MigrationIdFormat, str)
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
	Do() error
	Undo() error
	GetManager() Manager
	SetManager(manager Manager) Migration
}
