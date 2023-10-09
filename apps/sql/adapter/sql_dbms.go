package adapter

import (
	sqltypes "github.com/vertex-center/vertex/apps/sql/types"
)

type SqlDBMSAdapter struct{}

func NewSqlDBMSAdapter() sqltypes.SqlDBMSAdapterPort {
	return &SqlDBMSAdapter{}
}

func (a *SqlDBMSAdapter) GetDatabases() (*[]sqltypes.SqlDatabase, error) {
	// By default, return an empty list of databases. This should
	// be implemented by the specific adapter.
	return nil, nil
}
