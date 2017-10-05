package migration

type MigrationReporter interface {
	Error(err error)
	Warning(message string)
	Success(message string)
}
