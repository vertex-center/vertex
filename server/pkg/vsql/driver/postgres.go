package driver

var Postgres = PostgresDriver{}

type PostgresDriver struct{}

func (s PostgresDriver) AutoIncrement() string {
	return "SERIAL PRIMARY KEY"
}

func (s PostgresDriver) ForeignKey() string {
	return "FOREIGN KEY (%[1]s) REFERENCES %[2]s(%[3]s)"
}
