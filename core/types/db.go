package types

import "gorm.io/gorm"

type DB struct {
	*gorm.DB
}

func (d *DB) SetDB(db *gorm.DB) {
	d.DB = db
}
