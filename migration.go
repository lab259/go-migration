package migration

import "time"

const DEFAULT_MIGRATION_TABLE = "_migrations"

type Migration interface {
	GetId() time.Time
	GetDescription() string
	Up() error
	Down() error
	GetManager() Manager
	SetManager(manager Manager) Migration
}
