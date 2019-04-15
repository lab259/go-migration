package migration

import (
	"errors"
	"time"
)

// ErrMigrationPanicked is the error returned when a migrations panics.
var ErrMigrationPanicked = errors.New("migration panicked")

// ErrMigrationStarved is when late migrations are detected.
var ErrMigrationStarved = errors.New("migration starvation")

// ManagerDefault is a default implementation of a Manager. It provides, via
// migration.NewManager, a way to define what is the source and target of a
// manager.
type ManagerDefault struct {
	source Source
	target Target
}

// NewDefaultManager creates and returns a migration.Manager implementation
// (`migration.ManagerDefault`) based on a target and source.
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
	migrations, err := manager.source.List()
	if err != nil {
		return nil, err
	}

	executed, err := manager.target.MigrationsExecuted()
	if err != nil {
		return nil, err
	}
	executedMap := make(map[int64]bool, len(executed))
	for _, m := range executed {
		executedMap[m.UnixNano()] = true
	}

	pending := make([]Migration, 0, len(migrations))
	for _, migration := range migrations {
		if _, ok := executedMap[migration.GetID().UnixNano()]; !ok {
			pending = append(pending, migration)
		}
	}

	return pending, nil
}

// MigrationsExecuted returns a list of migrations that were executed. It uses
// the migration.Manager.MigrationsAfter passing on the current version from
// migration.Manager.Target.Version.
func (manager *ManagerDefault) MigrationsExecuted() ([]Migration, error) {
	migrations, err := manager.source.List()
	if err != nil {
		return nil, err
	}

	executed, err := manager.target.MigrationsExecuted()
	if err != nil {
		return nil, err
	}
	executedMap := make(map[int64]bool, len(executed))
	for _, m := range executed {
		executedMap[m.UnixNano()] = true
	}

	pending := make([]Migration, 0, len(migrations))
	for _, migration := range migrations {
		if _, ok := executedMap[migration.GetID().UnixNano()]; ok {
			pending = append(pending, migration)
		}
	}

	return pending, nil
}

// Do takes a step up on the migrations, bringing the database one step closer
// to the latest migration.
//
// Before the execution of the migrations, it calls the reporter.Before method.
// After the migration is executed, if it returns no error, it calls the
// reporter.After method.
func (manager *ManagerDefault) Do(reporter Reporter, executionContext interface{}) (*Summary, error) {
	version, err := manager.target.Version()
	if err != nil {
		return nil, err
	}
	migrations, err := manager.MigrationsPending()
	if err != nil {
		return nil, err
	}

	err = manager.detectStarvation(reporter, migrations, version)
	if err != nil {
		return nil, err
	}

	// No migrations
	if len(migrations) == 0 {
		return nil, nil
	}
	return manager.do(migrations[0], reporter, executionContext)
}

func (manager *ManagerDefault) do(m Migration, reporter Reporter, executionContext interface{}) (summary *Summary, err error) {
	summary = &Summary{
		Migration: m,
		direction: DirectionDo,
	}
	reporter.BeforeMigration(*summary, nil)

	startedAt := time.Now()
	func() {
		defer func() {
			if r := recover(); r != nil {
				summary.panicked = true
				if err, ok := r.(error); ok {
					summary.setFailed(err)
				}
				summary.panicData = r
				err = ErrMigrationPanicked
			}
		}()
		err = m.Do(executionContext)
	}()
	summary.duration = time.Since(startedAt)

	if !summary.panicked && err != nil {
		summary.setFailed(err)
		reporter.AfterMigration(*summary, err)
		return summary, err
	}
	reporter.AfterMigration(*summary, err)

	if summary.panicked {
		return summary, ErrMigrationPanicked
	}
	if err = manager.target.AddMigration(summary); err != nil {
		return summary, err
	}

	return summary, nil
}

// Undo takes a step up on the migrations, bringing the database one step closer
// to the latest migration.
//
// Before the execution of the migrations, it calls the reporter.Before method.
// After the migration is executed, if it returns no error, it calls the
// reporter.After method.
func (manager *ManagerDefault) Undo(reporter Reporter, executionContext interface{}) (*Summary, error) {
	migrations, err := manager.MigrationsExecuted()
	if err != nil {
		return nil, err
	}
	// No migrations
	if len(migrations) == 0 {
		return nil, nil
	}
	summary, err := manager.undo(migrations[len(migrations)-1], reporter, executionContext)
	return summary, err
}

func (manager *ManagerDefault) undo(m Migration, reporter Reporter, executionContext interface{}) (*Summary, error) {
	summary := &Summary{
		Migration: m,
		direction: DirectionUndo,
	}
	reporter.BeforeMigration(*summary, nil)

	startedAt := time.Now()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				summary.panicked = true
				// TODO Capture the panic info
				if err, ok := r.(error); ok {
					summary.setFailed(err)
				}
				summary.panicData = r
				err = ErrMigrationPanicked
			}
		}()
		err = m.Undo(executionContext)
	}()
	summary.duration = time.Since(startedAt)

	if !summary.panicked && err != nil {
		summary.setFailed(err)
		reporter.AfterMigration(*summary, err)
		return summary, err
	}
	reporter.AfterMigration(*summary, err)

	if summary.panicked {
		return summary, ErrMigrationPanicked
	}

	if err = manager.target.RemoveMigration(summary); err != nil {
		return summary, err
	}
	return summary, nil
}

func (manager *ManagerDefault) detectStarvation(reporter Reporter, list []Migration, version time.Time) error {
	migrationsStarved := make([]Migration, 0)

	for _, m := range list {
		if m.GetID().Before(version) {
			migrationsStarved = append(migrationsStarved, m)
		}
	}

	if len(migrationsStarved) > 0 {
		reporter.MigrationsStarved(migrationsStarved)
		return ErrMigrationStarved
	}
	return nil
}

// Migrate brings the database to the latest migration.
func (manager *ManagerDefault) Migrate(reporter Reporter, executionContext interface{}) ([]*Summary, error) {
	version, err := manager.target.Version()
	if err != nil {
		return nil, err
	}
	list, err := manager.MigrationsPending()
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		err = manager.detectStarvation(reporter, list, version)
		if err != nil {
			return nil, err
		}
	}
	reporter.BeforeMigrate(list)
	result := make([]*Summary, 0, len(list))
	for i := 0; i < len(list); i++ {
		summary, err := manager.do(list[i], reporter, executionContext)
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
func (manager *ManagerDefault) Rewind(reporter Reporter, executionContext interface{}) ([]*Summary, error) {
	list, err := manager.MigrationsExecuted()
	if err != nil {
		return nil, err
	}
	reporter.BeforeRewind(list)
	result := make([]*Summary, 0, len(list))
	for i := len(list) - 1; i > -1; i-- {
		summary, err := manager.undo(list[i], reporter, executionContext)
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
func (manager *ManagerDefault) Reset(reporter Reporter, executionContext interface{}) ([]*Summary, []*Summary, error) {
	reporter.BeforeReset()
	migrationsBack, err := manager.Rewind(reporter, executionContext)
	if err != nil {
		return migrationsBack, nil, err
	}
	migrationsForward, err := manager.Migrate(reporter, executionContext)
	return migrationsBack, migrationsForward, err
}
