package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/vsql"
)

var Migrations = []vsql.Migration{
	// v0.17.0
	v1{}, // add email to users
	v2{}, // remove timestamps from migrations
	v3{}, // remove email from users, add email table
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

type v3 struct{}

func (v3) Up(tx *sqlx.Tx) error {
	driver := vsql.DriverFromName(tx.DriverName())
	schema := vsql.BuildSchema(driver,
		vsql.AlterTable("users").
			RemoveField("email"),

		vsql.CreateTable("emails").
			WithID().
			WithField("user_id", "INTEGER", "NOT NULL").
			WithField("email", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithForeignKey("user_id", "users", "id"),
	)
	_, err := tx.Exec(schema)
	return err
}
