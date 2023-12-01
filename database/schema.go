package database

import (
	"time"

	"github.com/vertex-center/vertex/pkg/vsql"
)

func GetSchema(driver vsql.Driver) string {
	return vsql.BuildSchema(driver,
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

		vsql.CreateTable("migrations").
			WithID().
			WithField("version", "INTEGER", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt(),

		vsql.InsertInto("migrations").
			Columns("version", "created_at", "updated_at").
			Values(0, time.Now().Unix(), time.Now().Unix()),
	)
}
