package rlog

import (
	"fmt"
	"os"
	"time"

	rlog2 "github.com/lab259/rlog/v2"

	"github.com/lab259/go-migration"
)

const migrationIDFormat = "20060102150405"

// DefaultReporter is the default implementation of a Reporter.
type rlogReporter struct {
	exitFnc func(code int)
	logger  rlog2.Logger
}

func (reporter *rlogReporter) Failure(err error) {
	reporter.logger.Error(err)
}

// NewRLogReporter returns an instance of a DefaultReporter.
func NewRLogReporter(logger rlog2.Logger, exitFnc func(int)) *rlogReporter {
	return &rlogReporter{
		logger:  logger,
		exitFnc: exitFnc,
	}
}

// Exit calls the exit function of the reporter.
func (reporter *rlogReporter) Exit(code int) {
	reporter.exitFnc(code)
}

// BeforeMigration is called by the Manager right before a migration is ran.
func (reporter *rlogReporter) BeforeMigration(summary migration.Summary, err error) {
	var action string
	if summary.Direction() == migration.DirectionDo {
		action = "Applying"
	} else {
		action = "Rewinding"
	}
	reporter.logger.Tracef(2, "  %s [%s] %s ...", action, styleMigrationID(summary.Migration.GetID().Format(migrationIDFormat)), styleMigrationTitle(summary.Migration.GetDescription()))
}

// AfterMigration is called by the Manager right after a migrations is ran.
func (reporter *rlogReporter) AfterMigration(summary migration.Summary, err error) {
	var result string
	if summary.Failed() {
		result = styleError("Failed")
	} else if summary.Panicked() {
		result = styleError("Panicked")
	} else {
		result = styleSuccess("Ok")
	}

	d := summary.Duration()
	duration := styleSuccess(d.String())
	if d > time.Second*10 {
		duration = styleError(d.String())
	} else if d > time.Second*3 {
		duration = styleWarning(d.String())
	}

	if summary.Failed() {
		reporter.logger.Errorf("    %s (%s): %s", result, duration, summary.Failure())
	} else {
		reporter.logger.Tracef(2, "    %s (%s)", result, summary.Duration())
	}
}

// BeforeMigrate is called right before the process of migration is triggered.
func (reporter *rlogReporter) BeforeMigrate(migrations []migration.Migration) {
	if len(migrations) == 0 {
		reporter.noMigrationsPending()
	} else {
		reporter.logger.Tracef(1, "Preparing to apply %d migrations", len(migrations))
	}
}

// AfterMigrate is called right after the process of migration is completed.
func (reporter *rlogReporter) AfterMigrate(migrations []*migration.Summary, err error) {
	if err != nil {
		reporter.logger.Error("    %s", err)
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
	reporter.logger.Tracef(1, "%s migrations were applied %s failed", styleSuccess(executed), styleError(failed))
	if failed > 0 {
		reporter.Exit(10)
	}
}

// BeforeRewind is called right before the process of rewinding is triggered.
func (reporter *rlogReporter) BeforeRewind(migrations []migration.Migration) {
	if len(migrations) == 0 {
		reporter.noMigrationsExecuted()
	} else {
		reporter.logger.Tracef(1, fmt.Sprintf("Preparing to rewind %d migrations", len(migrations)))
	}
}

// AfterRewind is called right after the process of rewinding is completed.
func (reporter *rlogReporter) AfterRewind(migrations []*migration.Summary, err error) {
	executed := 0
	failed := 0
	for _, m := range migrations {
		if m.Failed() || m.Panicked() {
			failed++
		} else {
			executed++
		}
	}
	reporter.logger.Tracef(1, "%s migrations were rewinded %s failed", styleSuccess(executed), styleError(failed))
	if failed > 0 {
		reporter.Exit(10)
	}
}

// BeforeReset is called right before the process of reseting is triggered.
func (reporter *rlogReporter) BeforeReset() {
	// TODO
}

// AfterReset is called right before the process of reseting is completed.
func (reporter *rlogReporter) AfterReset(rewindSummary []*migration.Summary, migrateSummary []*migration.Summary, err error) {
	reporter.AfterRewind(rewindSummary, err)
	reporter.AfterMigrate(migrateSummary, err)
}

// MigrationSummary prints the summary of the migration.
func (reporter *rlogReporter) MigrationSummary(migration *migration.Summary, err error) {
	if migration == nil && err == nil {
		reporter.logger.Warn("Nothing to be done")
	}
}

func (reporter *rlogReporter) noMigrationsPending() {
	reporter.logger.Trace(2, "No migrations pending.")
}

// ListPending reports the list of the migrations peding to be executed.
func (reporter *rlogReporter) ListPending(migrations []migration.Migration, err error) {
	if err != nil {
		reporter.Failure(err)
		return
	}
	if len(migrations) == 0 {
		reporter.noMigrationsPending()
	} else {
		reporter.logger.Infof("%d migrations pending:", len(migrations))
		for i, m := range migrations {
			reporter.logger.Infof("%d) [%s] %s", i+1, styleMigrationID(m.GetID().Format(migrationIDFormat)), styleMigrationTitle(m.GetDescription()))
		}
	}
}

func (reporter *rlogReporter) noMigrationsExecuted() {
	reporter.logger.Info(styleWarning(fmt.Sprintf("No migrations were executed.")))
}

// ListExecuted reports the list of the migrations that were executed.
func (reporter *rlogReporter) ListExecuted(migrations []migration.Migration, err error) {
	if err != nil {
		reporter.Failure(err)
		return
	}
	if len(migrations) > 0 {
		reporter.logger.Infof("%d migrations executed:", len(migrations))
		for i, m := range migrations {
			reporter.logger.Infof("%d) [%s] %s", i+1, styleMigrationID(m.GetID().Format(migrationIDFormat)), styleMigrationTitle(m.GetDescription()))
		}
	} else {
		reporter.noMigrationsExecuted()
	}
}

// Usage prints the usage of the migration command.
func (reporter *rlogReporter) Usage() {
	reporter.logger.Info("Usage:", os.Args[0], "[migrate | rewind | do | undo | executed | pending]")
	line := "  %18s  %s"
	reporter.logger.Infof(line, styleBold("migrate"), "Apply all pending migrations")
	reporter.logger.Infof(line, styleBold("rewind"), "Rewind all executed migrations")
	reporter.logger.Infof(line, styleBold("do"), "Execute the next pending migration")
	reporter.logger.Infof(line, styleBold("undo"), "Execute the last applied migration")
	reporter.logger.Infof(line, styleBold("executed"), "List all executed migrations")
	reporter.logger.Infof(line, styleBold("pending"), "List all pending migrations")
}

// CommandNotFound reports the command executed by the migration tool was not
// found.
func (reporter *rlogReporter) CommandNotFound(command string) {
	reporter.logger.Errorf("Command %s was not found", command)
}

// NoCommand gets called when no command is provided to the migration tool.
func (reporter *rlogReporter) NoCommand() {
	reporter.Usage()
}

// MigrationsStarved is called whenever the manager detects a migration will
// never be executed.
func (reporter *rlogReporter) MigrationsStarved(migrations []migration.Migration) {
	reporter.logger.Error(styleError(fmt.Sprintf("Starvation detected in %d migrations", len(migrations))))
	for i, m := range migrations {
		reporter.logger.Infof("  %d) [%s] %s", i+1, styleMigrationID(m.GetID().Format(migrationIDFormat)), styleMigrationTitle(m.GetDescription()))
	}
}
