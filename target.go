package migration

import "time"

var NOVERSION time.Time = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

type Target interface {
	Version() (time.Time, error)
	SetVersion(id time.Time) error
}
