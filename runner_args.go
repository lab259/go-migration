package migration

import (
	"os"
)

// ArgsRunner is the Runner that will provide the default implementation of
// Runner that captures params from the arguments.
type ArgsRunner struct {
	reporter Reporter
	manager  Manager
	args     []string
	exitFnc  func(code int)
}

// NewArgsRunner returns a new instance of a ArgsRunner.
func NewArgsRunner(reporter Reporter, manager Manager, exitFnc func(code int)) Runner {
	return &ArgsRunner{
		reporter: reporter,
		manager:  manager,
		args:     os.Args[1:],
	}
}

// NewArgsRunnerCustom create a new instance of the Runer with custom arguments.
func NewArgsRunnerCustom(reporter Reporter, manager Manager, args ...string) Runner {
	return &ArgsRunner{
		reporter: reporter,
		manager:  manager,
		args:     args,
	}
}

// Run performs the actions based on the arguments captured.
func (runner *ArgsRunner) Run() {
	args := runner.args
	if len(args) > 0 {
		for _, s := range args {
			switch s {
			case "pending":
				runner.reporter.ListPending(runner.manager.MigrationsPending())
			case "executed":
				runner.reporter.ListExecuted(runner.manager.MigrationsExecuted())
			case "migrate":
				runner.reporter.AfterMigrate(runner.manager.Migrate(runner.reporter))
			case "rewind":
				runner.reporter.AfterRewind(runner.manager.Rewind(runner.reporter))
			case "do":
				runner.reporter.MigrationSummary(runner.manager.Do(runner.reporter))
			case "undo":
				runner.reporter.MigrationSummary(runner.manager.Undo(runner.reporter))
			case "reset":
				runner.reporter.AfterReset(runner.manager.Reset(runner.reporter))
			default:
				runner.reporter.CommandNotFound(s)
			}
			break
		}
	} else {
		runner.reporter.NoCommand()
	}
}
