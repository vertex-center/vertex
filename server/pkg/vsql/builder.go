package vsql

type Builder interface {
	Build(driver Driver) string
}
