package migration

// Reporter is a helper method used, mainly, for logging on the `Manager` methods.
// Before executing any migration (in any direction) the `Manager` calls the
// listener `Before`. Afterward, if it does not return any error, the listener's
// `After` is called.
type Reporter interface {
	BeforeMigration(migration MigrationSummary, err error)
	MigrationSummary(migration *MigrationSummary, err error)
	AfterMigration(migration MigrationSummary, err error)

	BeforeMigrate(migrations []Migration)
	AfterMigrate(migrations []*MigrationSummary, err error)

	BeforeRewind(migrations []Migration)
	AfterRewind(migrations []*MigrationSummary, err error)

	BeforeReset(doMigrations []Migration, undoMigrations []Migration)
	AfterReset(rewindSummary []*MigrationSummary, migrateSummary []*MigrationSummary, err error)

	ListPending(migrations []Migration, err error)
	ListExecuted(migrations []Migration, err error)
}
