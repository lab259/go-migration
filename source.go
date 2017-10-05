package migration

type Source interface {
	List() ([]Migration, error)
}
