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
	errDataConfigNotFound       = errors.New("config.yml doesn't exists or could not be found")
	errDataConfigFailedToRead   = errors.New("failed to read config.yml")
	errDataConfigFailedToDecode = errors.New("failed to decode config.yml")
)

// DataConfigFSAdapter is an adapter to configure how Vertex will store data.
type DataConfigFSAdapter struct {
	configDir string
	config    types.DataConfig
}

type DataConfigFSAdapterParams struct {
	configDir string
}

func NewDataConfigFSAdapter(params *DataConfigFSAdapterParams) port.DataConfigAdapter {
	if params == nil {
		params = &DataConfigFSAdapterParams{}
	}
	if params.configDir == "" {
		params.configDir = path.Join(storage.Path, "data")
	}

	err := os.MkdirAll(params.configDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.configDir),
		)
		os.Exit(1)
	}

	adapter := &DataConfigFSAdapter{
		configDir: params.configDir,
		config: types.DataConfig{
			DbmsName: types.DbNameSqlite,
		},
	}

	err = adapter.read()
	if errors.Is(err, errDataConfigFailedToDecode) {
		log.Error(err)
	}

	return adapter
}

func (a *DataConfigFSAdapter) GetDataConfig() types.DataConfig {
	return a.config
}

func (a *DataConfigFSAdapter) GetDBMSName() types.DbmsName {
	return a.config.DbmsName
}

func (a *DataConfigFSAdapter) SetDBMSName(name types.DbmsName) error {
	a.config.DbmsName = name
	return a.write()
}

func (a *DataConfigFSAdapter) read() error {
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

func (a *DataConfigFSAdapter) write() error {
	p := path.Join(a.configDir, "config.yml")

	data, err := yaml.Marshal(&a.config)
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, os.ModePerm)
}
