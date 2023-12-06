package database

import (
	"github.com/vertex-center/vertex/pkg/vsql"
)

func GetSchema(driver vsql.Driver) string {
	return vsql.BuildSchema(driver,
		vsql.CreateMigrationTable(Migrations),

		vsql.CreateTable("users").
			WithID().
			WithField("username", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt(),

		vsql.CreateTable("emails").
			WithID().
			WithField("user_id", "INTEGER", "NOT NULL").
			WithField("email", "VARCHAR(255)", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithForeignKey("user_id", "users", "id"),

		vsql.CreateTable("credentials_argon2").
			WithID().
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
			WithField("credential_id", "INTEGER", "NOT NULL").
			WithField("user_id", "INTEGER", "NOT NULL").
			WithPrimaryKey("credential_id", "user_id").
			WithForeignKey("credential_id", "credentials_argon2", "id").
			WithForeignKey("user_id", "users", "id"),

		vsql.CreateTable("sessions").
			WithID().
			WithField("token", "VARCHAR(255)", "NOT NULL").
			WithField("user_id", "INTEGER", "NOT NULL").
			WithCreatedAt().
			WithUpdatedAt().
			WithDeletedAt().
			WithForeignKey("user_id", "users", "id"),
	)
}
