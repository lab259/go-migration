package migration

// Source describes the common interface for source of migrations.
//
// Different technologies support various ways of interaction with them. For SQL
// based databases, SQL files may solve the problem. However, when dealing with
// NoSQL databases, as MongoDB, it will not work.
//
// So, migration.Source is a interface let the developer "store" migrations in
// many ways.
type Source interface {
	// List lists all migrations available for this migrations.Source.
	List() ([]Migration, error)
}
