package migration

import (
	"log"
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
	reporter.print(styleNormal("  Applying ["))
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
	reporter.printLn(fmt.Sprintf("Preparing to run %d migrations", len(migrations)))
}

func (reporter *DefaultReporter) AfterMigrate(migrations []*MigrationSummary, err error) {
	log.Println("AfterMigrate")
}

func (reporter *DefaultReporter) BeforeRewind(migrations []Migration) {
	log.Println("RewindSummary")
}

func (reporter *DefaultReporter) AfterRewind(migrations []*MigrationSummary, err error) {
	log.Println("RewindSummary")
}

func (reporter *DefaultReporter) BeforeReset(rewindSummary []Migration, migrateSummary []Migration) {
	log.Println("ResetSummary")
}

func (reporter *DefaultReporter) AfterReset(rewindSummary []*MigrationSummary, migrateSummary []*MigrationSummary, err error) {
	log.Println("ResetSummary")
}

func (reporter *DefaultReporter) MigrationSummary(migration *MigrationSummary, err error) {
	if migration == nil && err == nil {
		// TODO
		log.Println("Nothing to be done")
	} else if migration == nil {
		// TODO
		log.Println(err)
	} else if err != nil {
		// TODO
		log.Println(migration.Migration.GetDescription(), err)
	} else {
		// TODO
		log.Println(migration.Migration.GetDescription())
	}
}

func (reporter *DefaultReporter) noMigrationsPending() {
	reporter.printLn(styleSuccess("  No migrations pending."))
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
