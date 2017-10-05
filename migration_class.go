package migration

import "time"

// BaseMigration is the default implementation of the migration.Migration.
//
// It is designed to provide a coded implementaiton of a migration. It receives
// an up and down anonymous methods to be ran while executing the migration.
//
// This implementation is used by the migration.CodeSource implemenation of the
// migration.Source.
type BaseMigration struct {
	id          time.Time
	description string
	up          Handler
	down        Handler
	manager     Manager
}

// Handler is the signature of the up and down methods that a migration
// will receive.
type Handler func() error

// NewMigration returns a new instance of migration.Migration with all the
// required properties initialized.
//
// If a handler is provided it will assigned to the Up method. If a second is
// provided, it will be assigned to the Down method.
func NewMigration(id time.Time, description string, handlers ...Handler) *BaseMigration {
	var up, down Handler
	if len(handlers) > 0 {
		up = handlers[0]
	}
	if len(handlers) > 1 {
		down = handlers[1]
	}
	return &BaseMigration{
		id:          id,
		description: description,
		up:          up,
		down:        down,
	}
}

// GetID returns the ID of the migration.
func (m *BaseMigration) GetID() time.Time {
	return m.id
}

// GetDescription returns the ID of the migration.
func (m *BaseMigration) GetDescription() string {
	return m.description
}

// Up calls the up action of the migration.
func (m *BaseMigration) Up() error {
	return m.up()
}

// Down calls the down action of the migration.
func (m *BaseMigration) Down() error {
	return m.down()
}

// GetManager returns the reference of the manager that is executing the
// migration.
func (m *BaseMigration) GetManager() Manager {
	return m.manager
}

// SetManager set the reference of the manager that is executing the migration.
//
// It returns itself for sugar syntax.
func (m *BaseMigration) SetManager(manager Manager) Migration {
	m.manager = manager
	return m
}
