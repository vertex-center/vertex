package storage

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/vsql"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	errDbConfigNotFound       = errors.New("live/database/config.yml doesn't exists or could not be found")
	errDbConfigFailedToRead   = errors.New("failed to read live/database/config.yml")
	errDbConfigFailedToDecode = errors.New("failed to decode live/database/config.yml")
)

type DB struct {
	*sqlx.DB
	config DBConfig

	configPath string
}

type DBConfig struct {
	DbmsName   string `json:"dbms_name" yaml:"dbms_name"`
	DataSource string `json:"data_source" yaml:"data_source"`
}

type DBParams struct {
	// ID is the database identifier.
	// For an SQLite database, it is the path to the database file.
	// For a Postgres database, it is the database name.
	ID string

	// SchemaFunc is a function that returns the database schema depending on the driver.
	SchemaFunc func(driver vsql.Driver) string

	// Migrations is a list of migrations to run when needed.
	Migrations []vsql.Migration

	// ConfigDir is the path to the database config file.
	configPath string
}

func NewDB(params DBParams) (DB, error) {
	if params.ID == "" {
		params.ID = "default"
	}
	if params.configPath == "" {
		params.configPath = path.Join(FSPath, "database", "config.yml")
	}

	err := os.MkdirAll(path.Dir(params.configPath), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.configPath),
		)
		os.Exit(1)
	}

	db := DB{
		configPath: params.configPath,
		config: DBConfig{
			DbmsName:   "sqlite",
			DataSource: path.Join(FSPath, "database", params.ID+".db"),
		},
	}

	err = db.readConfig()
	if errors.Is(err, errDbConfigFailedToDecode) {
		return db, err
	}

	err = db.Connect()
	if err != nil {
		return db, err
	}

	err = db.runMigrations(params.SchemaFunc, params.Migrations)
	return db, err
}

func (db *DB) Connect() error {
	driver := db.config.DbmsName
	return db.ConnectTo(driver, db.config.DataSource, 10)
}

func (db *DB) ConnectTo(driver string, dataSource string, retries int) error {
	for i := 0; i < retries; i++ {
		conn, err := sqlx.Connect(driver, dataSource)
		if err != nil || db == nil {
			if i == retries-1 {
				return err
			}
			println("failed to connect to the database, retrying...")
			<-time.After(1 * time.Second)
		} else {
			db.DB = conn
			println("connected to the database after some retries")
			break
		}
	}
	return nil
}

func (db *DB) SetDBMSName(name string) error {
	db.config.DbmsName = name
	return db.writeConfig()
}

func (db *DB) GetDBMSName() string {
	return db.config.DbmsName
}

func (db *DB) readConfig() error {
	file, err := os.ReadFile(db.configPath)

	if errors.Is(err, fs.ErrNotExist) {
		return errDbConfigNotFound
	} else if err != nil {
		return errDbConfigFailedToRead
	}

	err = yaml.Unmarshal(file, &db.config)
	if err != nil {
		return errDbConfigFailedToDecode
	}
	return nil
}

func (db *DB) writeConfig() error {
	data, err := yaml.Marshal(&db.config)
	if err != nil {
		return err
	}
	return os.WriteFile(db.configPath, data, os.ModePerm)
}

func (db *DB) runMigrations(schemaFunc func(driver vsql.Driver) string, migrations []vsql.Migration) error {
	var current int
	err := db.Get(&current, "SELECT version FROM migrations LIMIT 1")
	if err != nil {
		return db.createSchemas(schemaFunc)
	}
	return vsql.Migrate(migrations, db.DB, current)
}

func (db *DB) createSchemas(schemaFunc func(driver vsql.Driver) string) error {
	vsqlDriver := vsql.DriverFromName(db.DriverName())
	schema := schemaFunc(vsqlDriver)
	_, err := db.Exec(schema)
	return err
}
