package migration

import "time"

// Summary is the record that keeps all the data collected while the migration
// was ran.
type Summary struct {
	Migration Migration
	direction Direction
	duration  time.Duration
	failed    bool
	failure   error
	panicked  bool
	panicData interface{}
}

// NewSummary creates a new summary based on a migration instance.
func NewSummary(migration Migration) *Summary {
	return &Summary{
		Migration: migration,
	}
}

// Failed returns if the migration failed.
func (summary *Summary) Failed() bool {
	return summary.failed
}

func (summary *Summary) setFailed(e error) {
	summary.failed = true
	summary.failure = e
}

// Failure returns the reason because the migration failed.
func (summary *Summary) Failure() error {
	return summary.failure
}

// Panicked is a flag that is returned when the tests panicks.
func (summary *Summary) Panicked() bool {
	return summary.panicked
}

// PanicData is the stores the data from the panic recovery function.
func (summary *Summary) PanicData() interface{} {
	return summary.panicData
}

// Direction is the direction the migrations ran.
func (summary *Summary) Direction() Direction {
	return summary.direction
}

// Duration is how long the migration took to run.
func (summary *Summary) Duration() time.Duration {
	return summary.duration
}
