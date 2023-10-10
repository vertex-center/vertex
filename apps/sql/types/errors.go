package types

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrCodeSQLDatabaseNotFound                  router.ErrCode = "sql_database_not_found"
	ErrCodeFailedToConfigureSQLDatabaseInstance router.ErrCode = "failed_to_configure_sql_database_instance"
)
