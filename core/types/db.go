package types

import "github.com/jmoiron/sqlx"

type DB struct {
	*sqlx.DB
}

func (d *DB) SetDB(db *sqlx.DB) {
	d.DB = db
}
