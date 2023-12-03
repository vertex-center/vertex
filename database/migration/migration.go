package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/vsql"
)

var Migrations = []vsql.Migration{
	// v0.17.0
	v1{}, // add email to users
	v2{}, // remove timestamps from migrations
}

type v1 struct{}

func (v1) Up(tx *sqlx.Tx) error {
	driver := vsql.DriverFromName(tx.DriverName())
	schema := vsql.BuildSchema(driver,
		vsql.AlterTable("users").
			AddField("email", "VARCHAR(255)"),
	)
	_, err := tx.Exec(schema)
	return err
}

type v2 struct{}

func (v2) Up(tx *sqlx.Tx) error {
	driver := vsql.DriverFromName(tx.DriverName())
	schema := vsql.BuildSchema(driver,
		vsql.AlterTable("migrations").
			RemoveField("updated_at").
			RemoveField("created_at").
			RemoveField("deleted_at"),
	)
	_, err := tx.Exec(schema)
	return err
}