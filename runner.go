package migration

// Runner defines the basic contract for a something that might run a migration.
type Runner interface {
	Run(executionContext interface{})
}
