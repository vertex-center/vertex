package adapter

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

var (
	errDbConfigNotFound       = errors.New("live/database/config.yml doesn't exists or could not be found")
	errDbConfigFailedToRead   = errors.New("failed to read live/database/config.yml")
	errDbConfigFailedToDecode = errors.New("failed to decode live/database/config.yml")
)

// DbAdapter is an adapter to configure how Vertex will store data.
type DbAdapter struct {
	configDir string
	config    types.DbConfig
	db        *types.DB
}

type DbAdapterParams struct {
	configDir string
}

func NewDbAdapter(params *DbAdapterParams) port.DbAdapter {
	if params == nil {
		params = &DbAdapterParams{}
	}
	if params.configDir == "" {
		params.configDir = path.Join(storage.Path, "database")
	}

	err := os.MkdirAll(params.configDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.configDir),
		)
		os.Exit(1)
	}

	adapter := &DbAdapter{
		configDir: params.configDir,
		config: types.DbConfig{
			DbmsName: types.DbmsNameSqlite,
		},
		db: &types.DB{},
	}

	err = adapter.read()
	if errors.Is(err, errDbConfigFailedToDecode) {
		log.Error(err)
	}

	return adapter
}

func (a *DbAdapter) Get() *types.DB {
	return a.db
}

func (a *DbAdapter) Connect() error {
	log.Info("connecting to the database", vlog.String("dbms", string(a.config.DbmsName)))

	var err error
	switch a.config.DbmsName {
	case types.DbmsNameSqlite:
		p := path.Join(a.configDir, "vertex.db")
		err = a.ConnectTo("sqlite3", p, 1)
	case types.DbmsNamePostgres:
		err = a.ConnectTo("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", 30)
	default:
		err = errors.New("invalid dbms name")
	}
	return err
}

func (a *DbAdapter) ConnectTo(driver string, dataSource string, retries int) error {
	for i := 0; i < retries; i++ {
		db, err := sqlx.Connect(driver, dataSource)
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
			a.db.SetDB(db)
			log.Info("connected to the database after some retries",
				vlog.Int("count", i+1),
			)
			break
		}
	}
	return nil
}

func (a *DbAdapter) GetDbConfig() types.DbConfig {
	return a.config
}

func (a *DbAdapter) GetDBMSName() types.DbmsName {
	return a.config.DbmsName
}

// SetDBMSName sets the database management system name.
// The user must also Connect to the database afterwords.
func (a *DbAdapter) SetDBMSName(name types.DbmsName) error {
	a.config.DbmsName = name
	return a.write()
}

func (a *DbAdapter) read() error {
	p := path.Join(a.configDir, "config.yml")
	file, err := os.ReadFile(p)

	if errors.Is(err, fs.ErrNotExist) {
		return errDbConfigNotFound
	} else if err != nil {
		return errDbConfigFailedToRead
	}

	err = yaml.Unmarshal(file, &a.config)
	if err != nil {
		return errDbConfigFailedToDecode
	}
	return nil
}

func (a *DbAdapter) write() error {
	p := path.Join(a.configDir, "config.yml")

	data, err := yaml.Marshal(&a.config)
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, os.ModePerm)
}
