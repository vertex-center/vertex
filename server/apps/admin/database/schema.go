package database

import (
	"time"

	"github.com/vertex-center/vertex/server/pkg/vsql"
)

func GetSchema(driver vsql.Driver) string {
	return vsql.BuildSchema(driver,
		vsql.CreateMigrationTable(Migrations),

		vsql.CreateTable("admin_settings").
			WithID().
			WithField("updates_channel", "VARCHAR(255)", "NOT NULL DEFAULT 'stable'").
			WithField("webhook", "VARCHAR(255)").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt(),

		vsql.InsertInto("admin_settings").
			Columns("updates_channel", "created_at", "updated_at").
			Values("stable", time.Now().Unix(), time.Now().Unix()),
	)
}
