package adapter

import (
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/types"
)

type SqlDBMSAdapter struct{}

func NewSqlDBMSAdapter() port.DBMSAdapter {
	return &SqlDBMSAdapter{}
}

func (a *SqlDBMSAdapter) GetDatabases() (*[]types.DB, error) {
	// By default, return an empty list of databases. This should
	// be implemented by the specific adapter.
	return nil, nil
}
