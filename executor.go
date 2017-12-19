package migration

import "bufio"

type Executor struct {
	reporter Reporter
	manager  Manager
	writer   *bufio.Writer
}

func (runner *Executor) Create() {
	// TODO
	panic("Not implemented")
}

func (runner *Executor) Migrate() {
	runner.reporter.AfterMigrate(runner.manager.Migrate(runner.reporter))
}

func (runner *Executor) Rewind() {
	runner.reporter.AfterRewind(runner.manager.Rewind(runner.reporter))
}

func (runner *Executor) Pending() {
	runner.reporter.ListPending(runner.manager.MigrationsPending())
}

func (runner *Executor) Executed() {
	runner.reporter.ListExecuted(runner.manager.MigrationsExecuted())
}

func (runner *Executor) Reset() {
	runner.reporter.AfterReset(runner.manager.Reset(runner.reporter))
}

func (runner *Executor) Do() {
	runner.reporter.MigrationSummary(runner.manager.Do(runner.reporter))
}

func (runner *Executor) Undo() {
	runner.reporter.MigrationSummary(runner.manager.Undo(runner.reporter))
}
