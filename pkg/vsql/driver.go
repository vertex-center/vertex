package vsql

import "github.com/vertex-center/vertex/pkg/vsql/driver"

type Driver interface {
	AutoIncrement() string
	ForeignKey() string
}

func DriverFromName(name string) Driver {
	switch name {
	case "postgres":
		return driver.Postgres
	default:
		return driver.Sqlite
	}
}
