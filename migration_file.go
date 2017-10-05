package migration

import "time"

type MigrationFile struct {
	id          time.Time
	description string
	baseFile    string
	ext         string
	up          bool
	down        bool
	manager     Manager
}

func (this *MigrationFile) GetId() time.Time {
	return this.id
}

func (this *MigrationFile) GetDescription() string {
	return this.description
}

func (this *MigrationFile) Up() error {
	return nil
}

func (this *MigrationFile) Down() error {
	return nil
}

func (this *MigrationFile) GetManager() Manager {
	return this.manager
}

func (this *MigrationFile) SetManager(manager Manager) Migration {
	this.manager = manager
	return this
}
