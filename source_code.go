package migration

// CodeSource is migration.Source implementation. It provides the development
// of migrations using Golang code.
//
// TODO add examples
type CodeSource struct {
	migrations []Migration
}

var defaultCode *CodeSource

func init() {
	defaultCode = NewCodeSource()
}

// DefaultCodeSource implements a singleton of the default CodeSource. It enable
// the developer to use the migration.Register without needing to deal with the
// migration.NewCodeSource.
//
// On normal situations only one CodeSource is enough through the whole system.
// And many sources are needed, the developer can use multiple instances of the
// migration.CodeSource and use the migration.CodeSource.Register, instead of
// using the default implementation migration.Register.
func DefaultCodeSource() *CodeSource {
	return defaultCode
}

// NewCodeSource returns a new instance of a migration.CodeSource.
func NewCodeSource() *CodeSource {
	return &CodeSource{
		migrations: make([]Migration, 0),
	}
}

// List implements the migration.Source.List by listing all the registered
// migrations of this instance.
func (s *CodeSource) List() ([]Migration, error) {
	return s.migrations, nil
}

// Register registers the migration for further use.
func (s *CodeSource) Register(migration Migration) {
	for i := 0; i < len(s.migrations); i++ {
		if s.migrations[i].GetID().After(migration.GetID()) {
			s.migrations = append(s.migrations[:i], append([]Migration{migration}, s.migrations[i:]...)...)
			return
		}
	}
	s.migrations = append(s.migrations, migration)
}

// Register registers the migration on the migration.DefaultCodeSource instance.
func Register(migration Migration) {
	defaultCode.Register(migration)
}
