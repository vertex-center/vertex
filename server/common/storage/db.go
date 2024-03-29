package storage

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/config"
	"github.com/vertex-center/vertex/server/pkg/vsql"
	"github.com/vertex-center/vlog"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sqlx.DB

	host, port, user, pass, driver, name string
}

type DBParams struct {
	Host, Port, User, Pass, Driver, Name string

	// SchemaFunc is a function that returns the database schema depending on the driver.
	SchemaFunc func(driver vsql.Driver) string

	// Migrations is a list of migrations to run when needed.
	Migrations []vsql.Migration
}

func NewDB(params DBParams) (DB, error) {
	driver, host, port, user, pass := config.Current.DB()
	if params.Driver == "" {
		params.Driver = driver
	}
	if params.Host == "" {
		params.Host = host
	}
	if params.Port == "" {
		params.Port = port
	}
	if params.User == "" {
		params.User = user
	}
	if params.Pass == "" {
		params.Pass = pass
	}
	if params.Name == "" {
		params.Name = "default"
	}

	db := DB{
		driver: params.Driver,
		host:   params.Host,
		port:   params.Port,
		user:   params.User,
		pass:   params.Pass,
		name:   params.Name,
	}

	err := db.Connect()
	if err != nil {
		return db, err
	}

	err = db.runMigrations(params.SchemaFunc, params.Migrations)
	return db, err
}

func (db *DB) Connect() error {
	switch db.driver {
	case "sqlite":
		dir := path.Join(FSPath, "db")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("%s.db", db.name)
		p := path.Join(dir, filename)
		return db.ConnectTo(db.driver, p, 1)
	case "postgres":
		source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db.host, db.port, db.user, db.pass, db.name)
		return db.ConnectTo(db.driver, source, 10)
	}
	return fmt.Errorf("unknown database driver %s", db.driver)
}

func (db *DB) ConnectTo(driver string, dataSource string, retries int) error {
	for i := 0; i < retries; i++ {
		conn, err := sqlx.Connect(driver, dataSource)
		if err != nil || db == nil {
			if i == retries-1 {
				return err
			}
			println("failed to connect to the database, retrying...")
			<-time.After(3 * time.Second)
		} else {
			db.DB = conn
			println("connected to the database after some retries")
			break
		}
	}
	return nil
}

func (db *DB) runMigrations(schemaFunc func(driver vsql.Driver) string, migrations []vsql.Migration) error {
	log.Info("running migrations for database", vlog.String("db_name", db.name))
	var current int
	err := db.Get(&current, "SELECT version FROM migrations LIMIT 1")
	if err != nil {
		return db.createSchemas(schemaFunc)
	}
	log.Info("database already initialized, running migrations instead", vlog.Int("current", current))
	return vsql.Migrate(migrations, db.DB, current)
}

func (db *DB) createSchemas(schemaFunc func(driver vsql.Driver) string) error {
	vsqlDriver := vsql.DriverFromName(db.DriverName())
	schema := schemaFunc(vsqlDriver)
	_, err := db.Exec(schema)
	return err
}
