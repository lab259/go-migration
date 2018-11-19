package migration

// Reporter is a helper method used, mainly, for logging on the `Manager` methods.
// Before executing any migration (in any direction) the `Manager` calls the
// listener `Before`. Afterward, if it does not return any error, the listener's
// `After` is called.
type Reporter interface {
	BeforeMigration(migration Summary, err error)
	MigrationSummary(migration *Summary, err error)
	AfterMigration(migration Summary, err error)

	BeforeMigrate(migrations []Migration)
	AfterMigrate(migrations []*Summary, err error)

	BeforeRewind(migrations []Migration)
	AfterRewind(migrations []*Summary, err error)

	BeforeReset(doMigrations []Migration, undoMigrations []Migration)
	AfterReset(rewindSummary []*Summary, migrateSummary []*Summary, err error)

	ListPending(migrations []Migration, err error)
	ListExecuted(migrations []Migration, err error)

	Failure(err error)
	Exit(code int)

	Usage()
	CommandNotFound(command string)
	NoCommand()
}
