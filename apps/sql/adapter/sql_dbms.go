package adapter

import (
	"github.com/vertex-center/vertex/apps/sql/types"
)

type SqlDBMSAdapter struct{}

func NewSqlDBMSAdapter() types.DBMSAdapterPort {
	return &SqlDBMSAdapter{}
}

func (a *SqlDBMSAdapter) GetDatabases() (*[]types.DB, error) {
	// By default, return an empty list of databases. This should
	// be implemented by the specific adapter.
	return nil, nil
}
