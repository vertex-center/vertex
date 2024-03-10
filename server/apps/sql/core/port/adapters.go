package port

import "github.com/vertex-center/vertex/server/apps/sql/core/types"

type DBMSAdapter interface {
	// GetDatabases returns a list of databases available in the DBMS.
	// If the DBMS does not support this operation, it will return nil.
	// If there is no database available, it will return an empty list.
	GetDatabases() (*[]types.DB, error)
}
