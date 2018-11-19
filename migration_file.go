package migration

import "time"

// FileMigration is the implementation of the migration.Migration that runs SQL
// files.
//
// It is designed to provide a coded implementaiton of a migration. It receives
// an up and down anonymous methods to be ran while executing the migration.
//
// It is used by the migration.CodeSource implemenation of the
// migration.Source.
type FileMigration struct {
	id          time.Time
	description string
	baseFile    string
	ext         string
	up          bool
	down        bool
	manager     Manager
}

// GetID implements the migration.Migration.GetID by returning the id of this
// migration.
func (m *FileMigration) GetID() time.Time {
	return m.id
}

// GetDescription implements the migration.Migration.GetDescription by returning the id of this
// migration.
func (m *FileMigration) GetDescription() string {
	return m.description
}

// Do implements the migration.Migration.Up by running all SQLs inside of the
// [migration.FileMigration.baseFile].up.sql file.
//
// If the file does not exists, it returns an error.
func (m *FileMigration) Do() error {
	// TODO
	return nil
}

// Undo implements the migration.Migration.Down by running all SQLs inside of
// the [migration.FileMigration.baseFile].down.sql file.
//
// If the file does not exists, it returns an error.
func (m *FileMigration) Undo() error {
	// TODO
	return nil
}

// GetManager implements the migration.Migration.GetManager by returning the
// manager responsible for the execution.
func (m *FileMigration) GetManager() Manager {
	return m.manager
}

// SetManager implements the migration.Migration.SetManager by setting the
// manager of this instance.
//
// It returns the itself for sugar syntax purposes.
func (m *FileMigration) SetManager(manager Manager) Migration {
	m.manager = manager
	return m
}
