package migration

import "time"

type MigrationBase struct {
	MigrationId time.Time
	description string
	up          MigrationHandler
	down        MigrationHandler
	manager     Manager
}

type MigrationHandler func() error

func NewMigration(id time.Time, description string, handlers ...MigrationHandler) *MigrationBase {
	var up, down MigrationHandler
	if len(handlers) > 0 {
		up = handlers[0]
	}
	if len(handlers) > 1 {
		down = handlers[1]
	}
	return &MigrationBase{
		MigrationId: id,
		description: description,
		up:          up,
		down:        down,
	}
}

func (this *MigrationBase) GetId() time.Time {
	return this.MigrationId
}

func (this *MigrationBase) GetDescription() string {
	return this.description
}

func (this *MigrationBase) Up() error {
	return this.up()
}

func (this *MigrationBase) Down() error {
	return this.down()
}

func (this *MigrationBase) GetManager() Manager {
	return this.manager
}

func (this *MigrationBase) SetManager(manager Manager) Migration {
	this.manager = manager
	return this
}
