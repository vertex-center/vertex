package driver

var Sqlite = SQLiteDriver{}

type SQLiteDriver struct{}

func (s SQLiteDriver) AutoIncrement() string {
	return "INTEGER PRIMARY KEY AUTOINCREMENT"
}

func (s SQLiteDriver) ForeignKey() string {
	return "FOREIGN KEY (%[1]s) REFERENCES %[2]s(%[3]s)"
}
