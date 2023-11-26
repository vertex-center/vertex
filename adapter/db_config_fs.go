package adapter

import (
	"errors"
	"io/fs"
	"os"
	"path"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

var (
	errDataConfigNotFound       = errors.New("live/database/config.yml doesn't exists or could not be found")
	errDataConfigFailedToRead   = errors.New("failed to read live/database/config.yml")
	errDataConfigFailedToDecode = errors.New("failed to decode live/database/config.yml")
)

// DbConfigFSAdapter is an adapter to configure how Vertex will store data.
type DbConfigFSAdapter struct {
	configDir string
	config    types.DbConfig
}

type DbConfigFSAdapterParams struct {
	configDir string
}

func NewDataConfigFSAdapter(params *DbConfigFSAdapterParams) port.DbConfigAdapter {
	if params == nil {
		params = &DbConfigFSAdapterParams{}
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

	adapter := &DbConfigFSAdapter{
		configDir: params.configDir,
		config: types.DbConfig{
			DbmsName: types.DbmsNameSqlite,
		},
	}

	err = adapter.read()
	if errors.Is(err, errDataConfigFailedToDecode) {
		log.Error(err)
	}

	return adapter
}

func (a *DbConfigFSAdapter) GetDbConfig() types.DbConfig {
	return a.config
}

func (a *DbConfigFSAdapter) GetDBMSName() types.DbmsName {
	return a.config.DbmsName
}

func (a *DbConfigFSAdapter) SetDBMSName(name types.DbmsName) error {
	a.config.DbmsName = name
	return a.write()
}

func (a *DbConfigFSAdapter) read() error {
	p := path.Join(a.configDir, "config.yml")
	file, err := os.ReadFile(p)

	if errors.Is(err, fs.ErrNotExist) {
		return errDataConfigNotFound
	} else if err != nil {
		return errDataConfigFailedToRead
	}

	err = yaml.Unmarshal(file, &a.config)
	if err != nil {
		return errDataConfigFailedToDecode
	}
	return nil
}

func (a *DbConfigFSAdapter) write() error {
	p := path.Join(a.configDir, "config.yml")

	data, err := yaml.Marshal(&a.config)
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, os.ModePerm)
}
