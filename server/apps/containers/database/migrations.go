package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/vsql"
)

var Migrations = []vsql.Migration{
	// 0.16
	&v1{}, // Rename service_id to template_id
	&v2{}, // Make template_id nullable
}

type v1 struct{}

func (m *v1) Up(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE containers
		RENAME COLUMN service_id TO template_id;
	`)
	return err
}

type v2 struct{}

func (m *v2) Up(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE containers
		ALTER COLUMN template_id DROP NOT NULL;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE containers
		SET template_id = NULL
		WHERE template_id = '';
	`)
	return err
}
