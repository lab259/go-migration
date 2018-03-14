package migration

import (
	"os"
	"fmt"
	"io"
)

type DefaultReporter struct {
	writer io.Writer
}

func NewDefaultReporter() *DefaultReporter {
	return &DefaultReporter{
		writer: os.Stdout,
	}
}

func (reporter *DefaultReporter) printLn(args ...interface{}) {
	fmt.Fprintln(reporter.writer, args...)
}

func (reporter *DefaultReporter) print(args ...interface{}) {
	fmt.Fprint(reporter.writer, args...)
}

func (reporter *DefaultReporter) BeforeMigration(summary MigrationSummary, err error) {
	if summary.Direction() == MigrationDirectionDo {
		reporter.print(styleNormal("  Applying ["))
	} else {
		reporter.print(styleNormal("  Rewinding ["))
	}
	reporter.print(styleMigrationId(summary.Migration.GetID().Format(MigrationIdFormat)))
	reporter.print(styleNormal("] "))
	reporter.print(styleMigrationTitle(summary.Migration.GetDescription()))
	reporter.print(styleNormal("... "))
}

func (reporter *DefaultReporter) AfterMigration(summary MigrationSummary, err error) {
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
		reporter.printLn(fmt.Sprintf("    %s", styleError(summary.Failure().Message)))
		reporter.printLn()
	}
}

func (reporter *DefaultReporter) BeforeMigrate(migrations []Migration) {
	if len(migrations) == 0 {
		reporter.noMigrationsPending()
	} else {
		reporter.printLn(fmt.Sprintf("Preparing to apply %d migrations", len(migrations)))
	}
}

func (reporter *DefaultReporter) AfterMigrate(migrations []*MigrationSummary, err error) {
	if err != nil {
		reporter.printLn(err)
		os.Exit(11)
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
		os.Exit(10)
	}
}

func (reporter *DefaultReporter) BeforeRewind(migrations []Migration) {
	if len(migrations) == 0 {
		reporter.noMigrationsExecuted()
	} else {
		reporter.printLn(fmt.Sprintf("Preparing to rewind %d migrations", len(migrations)))
	}
}

func (reporter *DefaultReporter) AfterRewind(migrations []*MigrationSummary, err error) {
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
		os.Exit(10)
	}
}

func (reporter *DefaultReporter) BeforeReset(rewindSummary []Migration, migrateSummary []Migration) {
	// TODO
}

func (reporter *DefaultReporter) AfterReset(rewindSummary []*MigrationSummary, migrateSummary []*MigrationSummary, err error) {
	reporter.printLn()
	reporter.AfterRewind(rewindSummary, err)
	reporter.AfterMigrate(migrateSummary, err)
	reporter.printLn()
}

func (reporter *DefaultReporter) MigrationSummary(migration *MigrationSummary, err error) {
	if migration == nil && err == nil {
		reporter.printLn(styleWarning("  Nothing to be done"))
	}
	reporter.printLn()
}

func (reporter *DefaultReporter) noMigrationsPending() {
	reporter.printLn(styleSuccess("  No migrations pending."))
	reporter.printLn()
}

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
			reporter.print(styleMigrationId(m.GetID().Format(MigrationIdFormat)))
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

func (reporter *DefaultReporter) ListExecuted(migrations []Migration, err error) {
	if err != nil {
		reporter.printLn(styleError(err.Error()))
		return
	}
	if len(migrations) > 0 {
		reporter.printLn(styleWarning(fmt.Sprintf("  %d migrations executed:", len(migrations))))
		for i, m := range migrations {
			reporter.print(styleNormal(fmt.Sprintf("  %d) [", i+1)))
			reporter.print(styleMigrationId(m.GetID().Format(MigrationIdFormat)))
			reporter.print(styleNormal("] "))
			reporter.printLn(styleMigrationTitle(m.GetDescription()))
		}
	} else {
		reporter.noMigrationsExecuted()
	}
	reporter.printLn()
}

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

func (reporter *DefaultReporter) CommandNotFound(command string) {
	reporter.printLn(fmt.Sprintf("Command %s was not found"))
}

func (reporter *DefaultReporter) NoCommand() {
	reporter.Usage()
}
