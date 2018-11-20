package migration

import (
	"fmt"
	"io"
	"os"
)

// DefaultReporter is the default implementation of a Reporter.
type DefaultReporter struct {
	writer  io.Writer
	exitFnc func(code int)
}

// NewDefaultReporter returns an instance of a DefaultReporter.
func NewDefaultReporter() *DefaultReporter {
	return NewDefaultReporterWithParams(os.Stdout, os.Exit)
}

// NewDefaultReporterWithParams initializes an instance of a DefaultReporter
// with its params.
func NewDefaultReporterWithParams(w io.Writer, exitFnc func(code int)) *DefaultReporter {
	return &DefaultReporter{
		writer:  w,
		exitFnc: exitFnc,
	}
}

func (reporter *DefaultReporter) printLn(args ...interface{}) {
	fmt.Fprintln(reporter.writer, args...)
}

func (reporter *DefaultReporter) print(args ...interface{}) {
	fmt.Fprint(reporter.writer, args...)
}

// Failure reports a failure and writes it down.
func (reporter *DefaultReporter) Failure(err error) {
	reporter.printLn(err)
}

// Exit calls the exit function of the reporter.
func (reporter *DefaultReporter) Exit(code int) {
	reporter.exitFnc(code)
}

// BeforeMigration is called by the Manager right before a migration is ran.
func (reporter *DefaultReporter) BeforeMigration(summary Summary, err error) {
	if summary.Direction() == DirectionDo {
		reporter.print(styleNormal("  Applying ["))
	} else {
		reporter.print(styleNormal("  Rewinding ["))
	}
	reporter.print(styleMigrationID(summary.Migration.GetID().Format(migrationIDFormat)))
	reporter.print(styleNormal("] "))
	reporter.print(styleMigrationTitle(summary.Migration.GetDescription()))
	reporter.print(styleNormal("... "))
}

// AfterMigration is called by the Manager right after a migrations is ran.
func (reporter *DefaultReporter) AfterMigration(summary Summary, err error) {
	if summary.Failed() {
		reporter.print(styleError("Failed"))
	} else if summary.Panicked() {
		reporter.print(styleError("Panicked"))
	} else {
		reporter.print(styleSuccess("Ok"))
	}
	ms := summary.duration.Nanoseconds() / 1000000
	if ms > 0 {
		if ms > 2000 {
			reporter.printLn(styleDurationVerySlow(fmt.Sprintf(" (%dms)", ms)))
		} else if ms > 500 {
			reporter.printLn(styleDurationSlow(fmt.Sprintf(" (%dms)", ms)))
		} else {
			reporter.printLn(styleDuration(fmt.Sprintf(" (%dms)", ms)))
		}
	} else {
		reporter.printLn()
	}
	if summary.Failed() {
		reporter.printLn()
		reporter.printLn(fmt.Sprintf("    %s", styleError(summary.Failure().Error())))
		reporter.printLn()
	}
}

// BeforeMigrate is called right before the process of migration is triggered.
func (reporter *DefaultReporter) BeforeMigrate(migrations []Migration) {
	if len(migrations) == 0 {
		reporter.noMigrationsPending()
	} else {
		reporter.printLn(fmt.Sprintf("Preparing to apply %d migrations", len(migrations)))
	}
}

// AfterMigrate is called right after the process of migration is completed.
func (reporter *DefaultReporter) AfterMigrate(migrations []*Summary, err error) {
	if err != nil {
		reporter.Failure(err)
		reporter.Exit(11)
	}
	executed := 0
	failed := 0
	for _, m := range migrations {
		if m.Failed() || m.Panicked() {
			failed++
		} else {
			executed++
		}
	}
	reporter.printLn(fmt.Sprintf("  %s migrations were applied %s failed", styleSuccess(fmt.Sprintf("%d", executed)), styleError(fmt.Sprintf("%d", failed))))
	if failed > 0 {
		reporter.Exit(10)
	}
}

// BeforeRewind is called right before the process of rewinding is triggered.
func (reporter *DefaultReporter) BeforeRewind(migrations []Migration) {
	if len(migrations) == 0 {
		reporter.noMigrationsExecuted()
	} else {
		reporter.printLn(fmt.Sprintf("Preparing to rewind %d migrations", len(migrations)))
	}
}

// AfterRewind is called right after the process of rewinding is completed.
func (reporter *DefaultReporter) AfterRewind(migrations []*Summary, err error) {
	executed := 0
	failed := 0
	for _, m := range migrations {
		if m.Failed() || m.Panicked() {
			failed++
		} else {
			executed++
		}
	}
	reporter.printLn(fmt.Sprintf("  %s migrations were rewinded %s failed", styleSuccess(fmt.Sprintf("%d", executed)), styleError(fmt.Sprintf("%d", failed))))
	if failed > 0 {
		reporter.Exit(10)
	}
}

// BeforeReset is called right before the process of reseting is triggered.
func (reporter *DefaultReporter) BeforeReset() {
	// TODO
}

// AfterReset is called right before the process of reseting is completed.
func (reporter *DefaultReporter) AfterReset(rewindSummary []*Summary, migrateSummary []*Summary, err error) {
	reporter.printLn()
	reporter.AfterRewind(rewindSummary, err)
	reporter.AfterMigrate(migrateSummary, err)
	reporter.printLn()
}

// MigrationSummary prints the summary of the migration.
func (reporter *DefaultReporter) MigrationSummary(migration *Summary, err error) {
	if migration == nil && err == nil {
		reporter.printLn(styleWarning("  Nothing to be done"))
	}
	reporter.printLn()
}

func (reporter *DefaultReporter) noMigrationsPending() {
	reporter.printLn(styleSuccess("  No migrations pending."))
	reporter.printLn()
}

// ListPending reports the list of the migrations peding to be executed.
func (reporter *DefaultReporter) ListPending(migrations []Migration, err error) {
	if err != nil {
		reporter.printLn(styleError(err.Error()))
		return
	}
	if len(migrations) == 0 {
		reporter.noMigrationsPending()
	} else {
		reporter.printLn(fmt.Sprintf("  %d migrations pending:", len(migrations)))
		for i, m := range migrations {
			reporter.print(styleNormal(fmt.Sprintf("    %d) [", i+1)))
			reporter.print(styleMigrationID(m.GetID().Format(migrationIDFormat)))
			reporter.print(styleNormal("] "))
			reporter.printLn(styleMigrationTitle(m.GetDescription()))
		}
	}
	reporter.printLn()
}

func (reporter *DefaultReporter) noMigrationsExecuted() {
	reporter.printLn(styleWarning(fmt.Sprintf("  No migrations were executed.")))
	reporter.printLn("")
}

// ListExecuted reports the list of the migrations that were executed.
func (reporter *DefaultReporter) ListExecuted(migrations []Migration, err error) {
	if err != nil {
		reporter.printLn(styleError(err.Error()))
		return
	}
	if len(migrations) > 0 {
		reporter.printLn(styleWarning(fmt.Sprintf("  %d migrations executed:", len(migrations))))
		for i, m := range migrations {
			reporter.print(styleNormal(fmt.Sprintf("  %d) [", i+1)))
			reporter.print(styleMigrationID(m.GetID().Format(migrationIDFormat)))
			reporter.print(styleNormal("] "))
			reporter.printLn(styleMigrationTitle(m.GetDescription()))
		}
	} else {
		reporter.noMigrationsExecuted()
	}
	reporter.printLn()
}

// Usage prints the usage of the migration command.
func (reporter *DefaultReporter) Usage() {
	reporter.printLn("Usage:", os.Args[0], "[migrate | rewind | do | undo | executed | pending]")
	reporter.printLn()
	line := "  %18s  %s"
	reporter.printLn(fmt.Sprintf(line, styleBold("migrate"), "Apply all pending migrations"))
	reporter.printLn(fmt.Sprintf(line, styleBold("rewind"), "Rewind all executed migrations"))
	reporter.printLn(fmt.Sprintf(line, styleBold("do"), "Execute the next pending migration"))
	reporter.printLn(fmt.Sprintf(line, styleBold("undo"), "Execute the last applied migration"))
	reporter.printLn(fmt.Sprintf(line, styleBold("executed"), "List all executed migrations"))
	reporter.printLn(fmt.Sprintf(line, styleBold("pending"), "List all pending migrations"))
	reporter.printLn()
}

// CommandNotFound reports the command executed by the migration tool was not
// found.
func (reporter *DefaultReporter) CommandNotFound(command string) {
	reporter.printLn(fmt.Sprintf("Command %s was not found", command))
}

// NoCommand gets called when no command is provided to the migration tool.
func (reporter *DefaultReporter) NoCommand() {
	reporter.Usage()
}

// MigrationsStarved is called whenever the manager detects a migration will
// never be executed.
func (reporter *DefaultReporter) MigrationsStarved(migrations []Migration) {
	reporter.printLn(styleError(fmt.Sprintf("Starvation detected in %d migrations", len(migrations))))
	for i, m := range migrations {
		reporter.print(styleNormal(fmt.Sprintf("  %d) [", i+1)))
		reporter.print(styleMigrationID(m.GetID().Format(migrationIDFormat)))
		reporter.print(styleNormal("] "))
		reporter.printLn(styleMigrationTitle(m.GetDescription()))
	}
}
