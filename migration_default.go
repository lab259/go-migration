package migration

import (
	"fmt"
	"path"
	"regexp"
	"runtime"
	"time"
)

type BaseMigration struct {
	id          time.Time
	description string
}

// GetID returns the ID of the migration.
func (m *BaseMigration) GetID() time.Time {
	return m.id
}

// GetDescription returns the ID of the migration.
func (m *BaseMigration) GetDescription() string {
	return m.description
}

// DefaultMigration is the default implementation of the migration.Migration.
//
// It is designed to provide a coded implementaiton of a migration. It receives
// an up and down anonymous methods to be ran while executing the migration.
//
// This implementation is used by the migration.CodeSource implemenation of the
// migration.Source.
type DefaultMigration struct {
	BaseMigration
	do      Handler
	undo    Handler
	manager Manager
}

// Handler is the signature of the up and down methods that a migration
// will receive.
type Handler func() error

// NewMigration returns a new instance of migration.Migration with all the
// required properties initialized.
//
// If a handler is provided it will assigned to the Up method. If a second is
// provided, it will be assigned to the Down method.
func NewMigration(id time.Time, description string, handlers ...Handler) *DefaultMigration {
	var do, undo Handler
	if len(handlers) > 0 {
		do = handlers[0]
	}
	if len(handlers) > 1 {
		undo = handlers[1]
	}
	return &DefaultMigration{
		BaseMigration: BaseMigration{
			id:          id,
			description: description,
		},
		do:   do,
		undo: undo,
	}
}

var codeMigrationRegex = regexp.MustCompile("^([0-9]{4}[0-9]{2}[0-9]{2}[0-9]{2}[0-9]{2}[0-9]{2})_(.*).go$")

// NewCodeMigration uses the regex to extract data NewMigration.
//
// If a handler is provided it will assigned to the Up method. If a second is
// provided, it will be assigned to the Down method.
func NewCodeMigration(handlers ...Handler) *DefaultMigration {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		groups := codeMigrationRegex.FindStringSubmatch(path.Base(file))
		if len(groups) == 3 {
			id, err := time.Parse("20060102150405", groups[1])
			if err != nil {
				panic(fmt.Sprintf("the file name '%s' has an invalid datetime"))
			}
			return NewMigration(id, groups[2], handlers...)
		}
		panic(fmt.Sprintf("the file name '%s' has an invalid format"))
	} else {
		panic(fmt.Sprintf("the file name '%s' has an invalid format"))
	}
}

// Up calls the up action of the migration.
func (m *DefaultMigration) Do() error {
	return m.do()
}

// Down calls the down action of the migration.
func (m *DefaultMigration) Undo() error {
	return m.undo()
}

// GetManager returns the reference of the manager that is executing the
// migration.
func (m *DefaultMigration) GetManager() Manager {
	return m.manager
}

// SetManager set the reference of the manager that is executing the migration.
//
// It returns itself for sugar syntax.
func (m *DefaultMigration) SetManager(manager Manager) Migration {
	m.manager = manager
	return m
}
