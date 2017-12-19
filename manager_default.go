package migration

import (
	"time"
)

// ManagerDefault is a default implementation of a Manager. It provides, via
// migration.NewManager, a way to define what is the source and target of a
// manager.
type ManagerDefault struct {
	source Source
	target Target
}

// NewManager creates and returns a migration.Manager implementation
// (migration.ManagerDefault) based on a target and source.
func NewDefaultManager(target Target, source Source) Manager {
	return &ManagerDefault{
		target: target,
		source: source,
	}
}

// Source returns the migration source used for this manager.
func (manager *ManagerDefault) Source() Source {
	return manager.source
}

// Target returns the migration target used for this manager.
func (manager *ManagerDefault) Target() Target {
	return manager.target
}

// MigrationsBefore returns all the migrations listed before the given `version`
// (exclusive).
func (manager *ManagerDefault) migrationsBefore(version time.Time) ([]Migration, error) {
	migrations, err := manager.source.List()
	if err == nil {
		til := 0
		for i, m := range migrations {
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
func (manager *ManagerDefault) migrationsAfter(version time.Time) ([]Migration, error) {
	migrations, err := manager.source.List()
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
func (manager *ManagerDefault) MigrationsPending() ([]Migration, error) {
	version, err := manager.target.Version()
	if err == nil {
		return manager.migrationsAfter(version)
	}
	return nil, err
}

// MigrationsDone returns a list of migrations that were executed. It uses the
// migration.Manager.MigrationsAfter passing on the current version from
// migration.Manager.Target.Version.
func (manager *ManagerDefault) MigrationsExecuted() ([]Migration, error) {
	version, err := manager.target.Version()
	if err == nil {
		return manager.migrationsBefore(version)
	}
	return nil, err
}

// Do takes a step up on the migrations, bringing the database one step closer
// to the latest migration.
//
// Before the execution of the migrations, it calls the listener.Before method.
// After the migration is executed, if it returns no error, it calls the
// listener.After method.
func (manager *ManagerDefault) Do(listener Reporter) (*MigrationSummary, error) {
	migrations, err := manager.MigrationsPending()
	if err != nil {
		return nil, err
	}
	// No migrations
	if len(migrations) == 0 {
		return nil, nil
	}
	summary := &MigrationSummary{
		Migration: migrations[0],
		direction: MigrationDirectionDo,
	}
	listener.BeforeMigration(*summary, nil)

	startedAt := time.Now()
	func() {
		defer func() {
			/*
			if r := recover(); r != nil {
				summary.panicked = true
				// TODO Capture the panic info
			}
			*/
		}()
		err = migrations[0].Do()
	}()
	summary.duration = time.Since(startedAt)

	if err != nil {
		summary.setFailed(err)
		listener.AfterMigration(*summary, err)
		return nil, err
	}
	if err = manager.target.SetVersion(migrations[0].GetID()); err == nil {
		listener.AfterMigration(*summary, nil)
		return summary, nil
	}
	summary.setFailed(err)
	return summary, err
}

// Undo takes a step up on the migrations, bringing the database one step closer
// to the latest migration.
//
// Before the execution of the migrations, it calls the listener.Before method.
// After the migration is executed, if it returns no error, it calls the
// listener.After method.
func (manager *ManagerDefault) Undo(listener Reporter) (*MigrationSummary, error) {
	migrations, err := manager.MigrationsExecuted()
	if err != nil {
		return nil, err
	}
	// No migrations
	if len(migrations) == 0 {
		return nil, nil
	}
	summary := &MigrationSummary{
		Migration: migrations[len(migrations)-1],
		direction: MigrationDirectionUndo,
	}
	listener.BeforeMigration(*summary, nil)

	startedAt := time.Now()
	func() {
		defer func() {
			/*
			if r := recover(); r != nil {
				summary.panicked = true
				// TODO Capture the panic info
			}
			*/
		}()
		err = migrations[0].Do()
	}()
	summary.duration = time.Since(startedAt)

	if err != nil {
		summary.setFailed(err)
		listener.AfterMigration(*summary, err)
		return nil, err
	}
	nversion := summary.Migration.GetID().Add(-time.Millisecond)
	if err = manager.target.SetVersion(nversion); err == nil {
		listener.AfterMigration(*summary, nil)
		return summary, nil
	}
	summary.setFailed(err)
	return summary, err
}

// Migrate brings the database to the latest migration.
func (manager *ManagerDefault) Migrate(reporter Reporter) ([]*MigrationSummary, error) {
	version, err := manager.target.Version()
	if err != nil {
		return nil, err
	}
	list, err := manager.migrationsAfter(version)
	if err != nil {
		return nil, err
	}
	reporter.BeforeMigrate(list)
	result := make([]*MigrationSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		summary, err := manager.Do(reporter)
		if summary != nil {
			result = append(result, summary)
		}
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// Rewind removes all migrations previously migrated.
//
// It lists all the executed migrations and executes their
// migration.Migrate.Down in a inverted order, virtually bringing the database
// to its original form.
func (manager *ManagerDefault) Rewind(listener Reporter) ([]*MigrationSummary, error) {
	version, err := manager.target.Version()
	if err != nil {
		return nil, err
	}
	list, err := manager.migrationsBefore(version)
	if err != nil {
		return nil, err
	}
	listener.BeforeRewind(list)
	result := make([]*MigrationSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		summary, err := manager.Undo(listener)
		if summary != nil {
			result = append(result, summary)
		}
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// Reset rewind all the migrations, then migrates to the latest.
func (manager *ManagerDefault) Reset(listener Reporter) ([]*MigrationSummary, []*MigrationSummary, error) {
	migrationsBack, err := manager.Rewind(listener)
	if err != nil {
		return migrationsBack, nil, err
	}
	migrationsForward, err := manager.Migrate(listener)
	return migrationsBack, migrationsForward, err
}
