package migration

import (
	"log"
)

type CodeSource struct {
	migrations []Migration
}

var defaultCode *CodeSource

func init() {
	defaultCode = NewCodeSource()
}

func DefaultCodeSource() *CodeSource {
	return defaultCode
}

func NewCodeSource() *CodeSource {
	return &CodeSource{
		migrations: make([]Migration, 0),
	}
}

func (this *CodeSource) List() ([]Migration, error) {
	return this.migrations, nil
}

func (this *CodeSource) Register(migration Migration) {
	for i := 0; i < len(this.migrations); i++ {
		if this.migrations[i].GetId().After(migration.GetId()) {
			this.migrations = append(this.migrations[:i], append([]Migration{migration}, this.migrations[i:]...)...)
			return
		}
	}
	this.migrations = append(this.migrations, migration)
}

func Register(migration Migration) {
	log.Println("Migration", migration.GetId(), "registered")
	defaultCode.Register(migration)
}
