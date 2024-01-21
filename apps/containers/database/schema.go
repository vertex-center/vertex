package database

import "github.com/vertex-center/vertex/pkg/vsql"

func GetSchema(driver vsql.Driver) string {
	return vsql.BuildSchema(driver,
		vsql.CreateMigrationTable(Migrations),

		vsql.CreateTable("containers").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("template_id", "VARCHAR(255)", "NOT NULL").
			WithField("user_id", "VARCHAR(36)", "NOT NULL").
			WithField("image", "VARCHAR(255)", "NOT NULL").
			WithField("image_tag", "VARCHAR(255)", "NOT NULL").
			WithField("status", "VARCHAR(255)", "NOT NULL").
			WithField("launch_on_startup", "BOOLEAN").
			WithField("name", "VARCHAR(255)", "NOT NULL").
			WithField("description", "VARCHAR(255)").
			WithField("color", "VARCHAR(7)").
			WithField("icon", "VARCHAR(255)").
			WithField("command", "VARCHAR(255)"),

		vsql.CreateTable("env_variables").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("type", "VARCHAR(255)", "NOT NULL").
			WithField("name", "VARCHAR(255)", "NOT NULL").
			WithField("display_name", "VARCHAR(255)", "NOT NULL").
			WithField("value", "VARCHAR(255)", "NOT NULL").
			WithField("default_value", "VARCHAR(255)").
			WithField("description", "VARCHAR(255)").
			WithField("secret", "BOOLEAN", "NOT NULL", "DEFAULT FALSE").
			WithForeignKey("container_id", "containers", "id"),

		vsql.CreateTable("capabilities").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("name", "VARCHAR(255)", "NOT NULL").
			WithForeignKey("container_id", "containers", "id"),

		vsql.CreateTable("ports").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("internal_port", "VARCHAR(255)", "NOT NULL").
			WithField("external_port", "VARCHAR(255)", "NOT NULL").
			WithForeignKey("container_id", "containers", "id"),

		vsql.CreateTable("volumes").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("type", "VARCHAR(255)", "NOT NULL").
			WithField("internal_path", "VARCHAR(255)", "NOT NULL").
			WithField("external_path", "VARCHAR(255)", "NOT NULL").
			WithForeignKey("container_id", "containers", "id"),

		vsql.CreateTable("sysctls").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("name", "VARCHAR(255)", "NOT NULL").
			WithField("value", "VARCHAR(255)", "NOT NULL").
			WithForeignKey("container_id", "containers", "id"),

		vsql.CreateTable("tags").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("user_id", "VARCHAR(36)", "NOT NULL").
			WithField("name", "VARCHAR(255)", "NOT NULL"),

		vsql.CreateTable("container_tags").
			WithField("container_id", "VARCHAR(36)", "NOT NULL").
			WithField("tag_id", "VARCHAR(36)", "NOT NULL").
			WithPrimaryKey("container_id", "tag_id").
			WithForeignKey("container_id", "containers", "id").
			WithForeignKey("tag_id", "tags", "id"),
	)
}
