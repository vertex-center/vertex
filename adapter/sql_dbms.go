package adapter

import "github.com/vertex-center/vertex/types"

type SqlDBMSAdapter struct{}

func NewSqlDBMSAdapter() types.SqlDBMSAdapterPort {
	return &SqlDBMSAdapter{}
}

func (a *SqlDBMSAdapter) GetDatabases() (*[]types.SqlDatabase, error) {
	// By default, return an empty list of databases. This should
	// be implemented by the specific adapter.
	return nil, nil
}
