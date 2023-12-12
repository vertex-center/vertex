package database

import "github.com/vertex-center/vertex/pkg/vsql"

func GetSchema(driver vsql.Driver) string {
	return vsql.BuildSchema(driver,
		vsql.CreateMigrationTable(Migrations),

		vsql.CreateTable("containers").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("service_id", "VARCHAR(255)", "NOT NULL").
			WithField("user_id", "VARCHAR(36)", "NOT NULL").
			WithField("image", "VARCHAR(255)", "NOT NULL").
			WithField("status", "VARCHAR(255)", "NOT NULL").
			WithField("launch_on_startup", "BOOLEAN").
			WithField("display_name", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithField("version", "VARCHAR(255)"),

		vsql.CreateTable("env_variables").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("key", "VARCHAR(255)", "NOT NULL").
			WithField("value", "VARCHAR(255)", "NOT NULL").
			WithPrimaryKey("container_uuid", "key").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithForeignKey("container_uuid", "containers", "uuid"),

		vsql.CreateTable("tags").
			WithField("container_uuid", "VARCHAR(36)", "NOT NULL").
			WithField("tag", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithPrimaryKey("container_uuid", "tag"),

		vsql.CreateTable("ports").
			WithField("container_uuid", "VARCHAR(36)", "NOT NULL").
			WithField("internal_port", "INTEGER", "NOT NULL").
			WithField("external_port", "INTEGER", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithPrimaryKey("container_uuid", "internal_port", "external_port"),

		vsql.CreateTable("volumes").
			WithField("container_uuid", "VARCHAR(36)", "NOT NULL").
			WithField("internal_path", "VARCHAR(255)", "NOT NULL").
			WithField("external_path", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithPrimaryKey("container_uuid", "internal_path", "external_path"),
	)
}
