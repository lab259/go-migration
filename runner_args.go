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
		exitFnc:  exitFnc,
	}
}

// NewArgsRunnerCustom create a new instance of the Runer with custom arguments.
func NewArgsRunnerCustom(reporter Reporter, manager Manager, exitFnc func(code int), args ...string) Runner {
	return &ArgsRunner{
		reporter: reporter,
		manager:  manager,
		exitFnc:  exitFnc,
		args:     args,
	}
}

// Run performs the actions based on the arguments captured.
func (runner *ArgsRunner) Run(executionContext interface{}) {
	args := runner.args
	if len(args) > 0 {
		for _, s := range args {
			switch s {
			case "pending":
				runner.reporter.ListPending(runner.manager.MigrationsPending())
			case "executed":
				runner.reporter.ListExecuted(runner.manager.MigrationsExecuted())
			case "migrate":
				runner.reporter.AfterMigrate(runner.manager.Migrate(runner.reporter, executionContext))
			case "rewind":
				runner.reporter.AfterRewind(runner.manager.Rewind(runner.reporter, executionContext))
			case "do":
				runner.reporter.MigrationSummary(runner.manager.Do(runner.reporter, executionContext))
			case "undo":
				runner.reporter.MigrationSummary(runner.manager.Undo(runner.reporter, executionContext))
			case "reset":
				runner.reporter.AfterReset(runner.manager.Reset(runner.reporter, executionContext))
			default:
				runner.reporter.CommandNotFound(s)
			}
			break
		}
	} else {
		runner.reporter.NoCommand()
	}
}
