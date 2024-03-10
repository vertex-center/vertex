package adapter

import (
	"github.com/vertex-center/vertex/server/apps/sql/core/port"
	"github.com/vertex-center/vertex/server/apps/sql/core/types"
)

type sqlDBMSAdapter struct{}

func NewSqlDBMSAdapter() port.DBMSAdapter {
	return &sqlDBMSAdapter{}
}

func (a *sqlDBMSAdapter) GetDatabases() (*[]types.DB, error) {
	// By default, return an empty list of databases. This should
	// be implemented by the specific adapter.
	return nil, nil
}
