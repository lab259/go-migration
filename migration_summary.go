package migration

import "time"

type MigrationFailure struct {
	// TODO
	Message string
}

type MigrationSummary struct {
	Migration Migration
	direction MigrationDirection
	duration  time.Duration
	failed    bool
	failure   error
	panicked  bool
	panicData interface{}
}

func NewMigrationSummary(migration Migration) *MigrationSummary {
	return &MigrationSummary{
		Migration: migration,
	}
}

func (summary *MigrationSummary) Failed() bool {
	return summary.failed
}

func (summary *MigrationSummary) setFailed(e error) {
	summary.failed = true
	summary.failure = e
}

func (summary *MigrationSummary) Failure() error {
	return summary.failure
}

func (summary *MigrationSummary) Panicked() bool {
	return summary.panicked
}

func (summary *MigrationSummary) PanicData() interface{} {
	return summary.panicData
}

func (summary *MigrationSummary) Direction() MigrationDirection {
	return summary.direction
}
