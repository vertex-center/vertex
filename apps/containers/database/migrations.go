package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/vsql"
)

var Migrations = []vsql.Migration{
	// add migrations here
	&v1{}, // Rename ServiceID to TemplateID
}

type v1 struct{}

func (m *v1) Up(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE containers
		RENAME COLUMN service_id TO template_id;
	`)
	return err
}
