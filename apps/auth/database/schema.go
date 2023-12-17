package database

import (
	"github.com/vertex-center/vertex/pkg/vsql"
)

func GetSchema(driver vsql.Driver) string {
	return vsql.BuildSchema(driver,
		vsql.CreateMigrationTable(Migrations),

		vsql.CreateTable("users").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("username", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt(),

		vsql.CreateTable("emails").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("user_id", "VARCHAR(36)", "NOT NULL").
			WithField("email", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithForeignKey("user_id", "users", "id"),

		vsql.CreateTable("credentials_argon2").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("login", "VARCHAR(255)", "NOT NULL").
			WithField("hash", "VARCHAR(255)", "NOT NULL").
			WithField("type", "VARCHAR(255)", "NOT NULL").
			WithField("iterations", "INTEGER", "NOT NULL").
			WithField("memory", "INTEGER", "NOT NULL").
			WithField("parallelism", "INTEGER", "NOT NULL").
			WithField("salt", "VARCHAR(255)", "NOT NULL").
			WithField("key_len", "INTEGER", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt(),

		vsql.CreateTable("credentials_argon2_users").
			WithField("credential_id", "VARCHAR(36)", "NOT NULL").
			WithField("user_id", "VARCHAR(36)", "NOT NULL").
			WithPrimaryKey("credential_id", "user_id").
			WithForeignKey("credential_id", "credentials_argon2", "id").
			WithForeignKey("user_id", "users", "id"),

		vsql.CreateTable("sessions").
			WithField("id", "VARCHAR(36)", "NOT NULL", "PRIMARY KEY").
			WithField("token", "VARCHAR(255)", "NOT NULL").
			WithField("user_id", "VARCHAR(36)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithForeignKey("user_id", "users", "id"),
	)
}
