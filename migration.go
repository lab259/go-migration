package migration

import "time"

// DefaultMigrationTable is the default name of the migrations table.
const DefaultMigrationTable = "_migrations"

// Migration is the interface that describes the common behavior that a
// migration should have to be manageable by the migration.Manager.
type Migration interface {
	GetID() time.Time
	GetDescription() string
	Up() error
	Down() error
	GetManager() Manager
	SetManager(manager Manager) Migration
}
