package migration

import (
	"time"
)

type listenerHandler func(migration Migration)

// Listener is a helper method used, mainly, for logging on the `Manager` methods.
// Before executing any migration (in any direction) the `Manager` calls the
// listener `Before`. Afterward, if it does not return any error, the listener's
// `After` is called.
type Listener struct {
	Before listenerHandler
	After  listenerHandler
}

// Manager is an interface that describe the common behavior of a migration
// manager.
//
// Any manager is divider in two parts: the `Source` and `Target`. In short, the
// `Source` is the origin of the migrations (eg. SQL files or Go scripts). A
// `Target` is database technology that the migrations will be action on.
//
// Integrating the `Source` and the `Target`, the `Manager` is responsible for
// running migrations with its methods `Migrate`, `Rewind`, `Reset`, `Up` and
// `Down`.
type Manager interface {
	Source() Source
	Target() Target
	Migrate(*Listener) ([]Migration, error)
	MigrationsDone() ([]Migration, error)
	MigrationsPending() ([]Migration, error)
	Rewind(*Listener) ([]Migration, error)
	Reset(*Listener, *Listener) ([]Migration, []Migration, error)
	Up(listener *Listener) (Migration, error)
	Down(listener *Listener) (Migration, error)
}

// ManagerBase is a default implementation of a Manager. It provides, via
// migration.NewManager, a way to define what is the source and target of a
// manager.
type ManagerBase struct {
	source Source
	target Target
}

// NewManager creates and returns a migration.Manager implementation
// (migration.ManagerBase) based on a target and source.
func NewManager(target Target, source Source) Manager {
	return &ManagerBase{
		target: target,
		source: source,
	}
}

// Source returns the migration source used for this manager.
func (m *ManagerBase) Source() Source {
	return m.source
}

// Target returns the migration target used for this manager.
func (m *ManagerBase) Target() Target {
	return m.target
}

// MigrationsBefore returns all the migrations listed before the given `version`
// (exclusive).
func (m *ManagerBase) MigrationsBefore(version time.Time) ([]Migration, error) {
	migrations, err := m.source.List()
	if err == nil {
		til := 0
		for i := 0; i < len(migrations); i++ {
			m := migrations[i]
			if m.GetID().After(version) {
				return migrations[:i], nil
			}
			til = i + 1
		}
		return migrations[0:til], nil
	}
	return nil, err
}

// MigrationsAfter returns all the migrations listed after the given `version`
// (exclusive).
func (m *ManagerBase) MigrationsAfter(version time.Time) ([]Migration, error) {
	migrations, err := m.source.List()
	if err == nil {
		for i := 0; i < len(migrations); i++ {
			m := migrations[i]
			if m.GetID().After(version) {
				return migrations[i:], nil
			}
		}
		return migrations[0:0], nil
	}
	return nil, err
}

// MigrationsPending returns a list of migrations that were not executed yet. It
// uses the migration.Manager.MigrationsBefore passing on the current version
// from migration.Manager.Target.Version.
func (m *ManagerBase) MigrationsPending() ([]Migration, error) {
	version, err := m.target.Version()
	if err == nil {
		return m.MigrationsAfter(version)
	}
	return nil, err
}

// MigrationsDone returns a list of migrations that were executed. It uses the
// migration.Manager.MigrationsAfter passing on the current version from
// migration.Manager.Target.Version.
func (m *ManagerBase) MigrationsDone() ([]Migration, error) {
	version, err := m.target.Version()
	if err == nil {
		return m.MigrationsBefore(version)
	}
	return nil, err
}

// Up takes a step up on the migrations, bringing the database one step closer
// to the latest migration.
//
// Before the execution of the migrations, it calls the listener.Before method.
// After the migration is executed, if it returns no error, it calls the
// listener.After method.
func (m *ManagerBase) Up(listener *Listener) (Migration, error) {
	migrations, err := m.MigrationsPending()
	if err == nil {
		if len(migrations) > 0 {
			listener.Before(migrations[0])
			err = migrations[0].Up()
			if err == nil {
				if err = m.target.SetVersion(migrations[0].GetID()); err == nil {
					listener.After(migrations[0])
					return migrations[0], err
				}
				return nil, err
			}
			return nil, err
		}
		return nil, nil
	}
	return nil, err
}

// Down takes a step up on the migrations, bringing the database one step closer
// to the latest migration.
//
// Before the execution of the migrations, it calls the listener.Before method.
// After the migration is executed, if it returns no error, it calls the
// listener.After method.
func (m *ManagerBase) Down(listener *Listener) (Migration, error) {
	migrations, err := m.MigrationsDone()
	if err == nil {
		if len(migrations) > 0 {
			i := len(migrations) - 1
			listener.Before(migrations[i])
			err = migrations[i].Down()
			if err == nil {
				if err = m.target.SetVersion(migrations[i].GetID()); err == nil {
					listener.After(migrations[i])
					return migrations[i], err
				}
				return nil, err
			}
			return nil, err
		}
		return nil, nil
	}
	return nil, err
}

// Migrate brings the database to the latest migration.
func (m *ManagerBase) Migrate(listener *Listener) ([]Migration, error) {
	list, err := m.MigrationsPending()

	if err != nil {
		return nil, err
	}
	result := make([]Migration, 0, len(list))
	for i := 0; i < len(list); i++ {
		if (listener != nil) && (listener.Before != nil) {
			listener.Before(list[i])
		}
		if err = list[i].SetManager(m).Up(); err != nil {
			return result, err
		}
		result = append(result, list[i])
		if err := m.Target().SetVersion(list[i].GetID()); err != nil {
			return result, err
		}
		if (listener != nil) && (listener.After != nil) {
			listener.After(list[i])
		}
	}
	return result, nil
}

// Rewind removes all migrations previously migrated.
//
// It lists all the executed migrations and executes their
// migration.Migrate.Down in a inverted order, virtually bringing the database
// to its original form.
func (m *ManagerBase) Rewind(listener *Listener) ([]Migration, error) {
	version, err := m.target.Version()
	if err != nil {
		return nil, err
	}
	list, err := m.MigrationsBefore(version)
	if err != nil {
		return nil, err
	}
	result := make([]Migration, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		m.target.SetVersion(list[i].GetID())
		if (listener != nil) && (listener.Before != nil) {
			listener.Before(list[i])
		}
		if err := list[i].SetManager(m).Down(); err != nil {
			return result, err
		}
		if (listener != nil) && (listener.After != nil) {
			listener.After(list[i])
		}
		result = append(result, list[i])
	}
	m.target.SetVersion(NoVersion)
	return result, nil
}

// Reset rewind all the migrations, then migrates to the latest.
func (m *ManagerBase) Reset(listenerRewind *Listener, listenerMigrate *Listener) ([]Migration, []Migration, error) {
	migrationsBack, err := m.Rewind(listenerRewind)
	if err != nil {
		return migrationsBack, nil, err
	}
	migrationsForward, err := m.Migrate(listenerMigrate)
	if err != nil {
		return migrationsBack, migrationsForward, err
	}
	return migrationsBack, migrationsForward, nil
}
