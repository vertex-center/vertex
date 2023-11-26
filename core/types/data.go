package types

type DbmsName string

const (
	DbmsNameSqlite   DbmsName = "sqlite" // Default
	DbmsNamePostgres DbmsName = "postgres"
)

type DbConfig struct {
	// DbmsName is the database management system name that Vertex will use.
	DbmsName DbmsName `json:"dbms_name" yaml:"dbms_name"`
}
