package migration

import (
	"time"
)

type MigrationListenerHandler func(migration Migration)

type MigrationListener struct {
	Before MigrationListenerHandler
	After  MigrationListenerHandler
}
type Manager interface {
	Source() Source
	Target() Target
	Migrate(*MigrationListener) ([]Migration, error)
	MigrationsDone() ([]Migration, error)
	MigrationsPending() ([]Migration, error)
	Rewind(*MigrationListener) ([]Migration, error)
	Reset(*MigrationListener, *MigrationListener) ([]Migration, []Migration, error)
	Up(listener *MigrationListener) (Migration, error)
	Down(listener *MigrationListener) (Migration, error)
}

type ManagerBase struct {
	source Source
	target Target
}

func NewManager(target Target, source Source) Manager {
	return &ManagerBase{
		target: target,
		source: source,
	}
}

func (this *ManagerBase) Source() Source {
	return this.source
}

func (this *ManagerBase) Target() Target {
	return this.target
}

// This method returns all the migrations listed before the given `id`
// (exclusive).
func (this *ManagerBase) MigrationsBefore(id time.Time) ([]Migration, error) {
	if migrations, err := this.source.List(); err == nil {
		til := 0
		for i := 0; i < len(migrations); i++ {
			m := migrations[i]
			if m.GetId().After(id) {
				return migrations[:i], nil
			} else {
				til = i + 1
			}
		}
		return migrations[0:til], nil
	} else {
		return nil, err
	}
}

// This method returns all the migrations listed after the given `id`
// (exclusive).
func (this *ManagerBase) MigrationsAfter(id time.Time) ([]Migration, error) {
	if migrations, err := this.source.List(); err == nil {
		for i := 0; i < len(migrations); i++ {
			m := migrations[i]
			if m.GetId().After(id) {
				return migrations[i:], nil
			}
		}
		return migrations[0:0], nil
	} else {
		return nil, err
	}
}

func (this *ManagerBase) MigrationsPending() ([]Migration, error) {
	if version, err := this.target.Version(); err == nil {
		return this.MigrationsAfter(version)
	} else {
		return nil, err
	}
}

func (this *ManagerBase) MigrationsDone() ([]Migration, error) {
	if version, err := this.target.Version(); err == nil {
		return this.MigrationsBefore(version)
	} else {
		return nil, err
	}
}

// This method take a step up on the migrations, bringing the database one step
// closer to the latest migration.
func (this *ManagerBase) Up(listener *MigrationListener) (Migration, error) {
	if migrations, err := this.MigrationsPending(); err == nil {
		if len(migrations) > 0 {
			listener.Before(migrations[0])
			err = migrations[0].Up()
			if err == nil {
				if err = this.target.SetVersion(migrations[0].GetId()); err == nil {
					listener.After(migrations[0])
					return migrations[0], err
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

// This method take a step down on the migrations, bringing the database one
// step away to the latest migration.
func (this *ManagerBase) Down(listener *MigrationListener) (Migration, error) {
	if migrations, err := this.MigrationsDone(); err == nil {
		if len(migrations) > 0 {
			i := len(migrations) - 1
			listener.Before(migrations[i])
			err = migrations[i].Down()
			if err == nil {
				if err = this.target.SetVersion(migrations[i].GetId()); err == nil {
					listener.After(migrations[i])
					return migrations[i], err
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

// This method takes the database to the latest migration.
func (this *ManagerBase) Migrate(listener *MigrationListener) ([]Migration, error) {
	version, err := this.target.Version()
	if err != nil {
		return nil, err
	}
	list, err := this.MigrationsAfter(version)
	if err != nil {
		return nil, err
	}
	result := make([]Migration, 0, len(list))
	for i := 0; i < len(list); i++ {
		if (listener != nil) && (listener.Before != nil) {
			listener.Before(list[i])
		}
		if err := list[i].SetManager(this).Up(); err == nil {
			result = append(result, list[i])
			if err := this.Target().SetVersion(list[i].GetId()); err != nil {
				return result, err
			}
			if (listener != nil) && (listener.After != nil) {
				listener.After(list[i])
			}
		} else {
			return result, err
		}
	}
	return result, nil
}

// Reverts all the migrations
func (this *ManagerBase) Rewind(listener *MigrationListener) ([]Migration, error) {
	version, err := this.target.Version()
	if err != nil {
		return nil, err
	}
	list, err := this.MigrationsBefore(version)
	if err != nil {
		return nil, err
	}
	result := make([]Migration, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		this.target.SetVersion(list[i].GetId())
		if (listener != nil) && (listener.Before != nil) {
			listener.Before(list[i])
		}
		if err := list[i].SetManager(this).Down(); err != nil {
			return result, err
		}
		if (listener != nil) && (listener.After != nil) {
			listener.After(list[i])
		}
		result = append(result, list[i])
	}
	this.target.SetVersion(NOVERSION)
	return result, nil
}

// Rewind all the migrations, then migrates to the latest.
func (this *ManagerBase) Reset(listenerRewind *MigrationListener, listenerMigrate *MigrationListener) ([]Migration, []Migration, error) {
	if migrationsBack, err := this.Rewind(listenerRewind); err != nil {
		return migrationsBack, nil, err
	} else {
		if migrationsForward, err := this.Migrate(listenerMigrate); err != nil {
			return migrationsBack, migrationsForward, err
		} else {
			return migrationsBack, migrationsForward, nil
		}
	}
}
