package storage

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
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
	configDir string
}

func NewDB(params *DBParams) (DB, error) {
	if params == nil {
		params = &DBParams{}
	}
	if params.configDir == "" {
		params.configDir = path.Join(FSPath, "database", "config.yml")
	}

	db := DB{
		configPath: params.configDir,
		config: DBConfig{
			DbmsName:   "sqlite",
			DataSource: path.Join(FSPath, "database", "default.db"),
		},
	}

	err := db.readConfig()
	if errors.Is(err, errDbConfigFailedToDecode) {
		return db, err
	}

	err = db.Connect()
	if err != nil {
		return db, err
	}

	return db, nil
}

func (db *DB) Connect() error {
	driver := db.config.DbmsName
	log.Info("connecting to the database", vlog.String("driver", driver))
	return db.ConnectTo(driver, db.config.DataSource, 10)
}

func (db *DB) ConnectTo(driver string, dataSource string, retries int) error {
	for i := 0; i < retries; i++ {
		conn, err := sqlx.Connect(driver, dataSource)
		if err != nil || db == nil {
			if i == retries-1 {
				return err
			}
			log.Info("failed to connect to the database, retrying...",
				vlog.String("error", err.Error()),
				vlog.Int("retry", i+1),
			)
			<-time.After(1 * time.Second)
		} else {
			db.DB = conn
			log.Info("connected to the database after some retries", vlog.Int("count", i+1))
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
