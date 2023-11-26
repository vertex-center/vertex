package types

type DbmsName string

const (
	DbNameSqlite   DbmsName = "sqlite" // Default
	DbNamePostgres DbmsName = "postgres"
)

type DataConfig struct {
	// DbmsName is the database management system name that Vertex will use.
	DbmsName DbmsName `json:"dbms_name"`
}
