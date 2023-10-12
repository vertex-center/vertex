package migration

import (
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

type MigrationTool struct {
	livePath     string
	metadataPath string
	migrations   []Migration
}

type LiveVersion struct {
	Version int
}

type Migration interface {
	Up(livePath string) error
}

func NewMigrationTool(livePath string) *MigrationTool {
	return &MigrationTool{
		livePath:     livePath,
		metadataPath: path.Join(livePath, "metadata.yml"),
		migrations: []Migration{
			&migration0{},
		},
	}
}

func (t *MigrationTool) Migrate() error {
	v, err := t.readLiveVersion()
	if err != nil {
		return err
	}

	for i := v.Version + 1; i < len(t.migrations); i++ {
		log.Info("running migration", vlog.Int("version", i))
		err := t.migrations[i].Up(t.livePath)
		if err != nil {
			return err
		}

		v.Version = i
		err = t.writeLiveVersion(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *MigrationTool) readLiveVersion() (*LiveVersion, error) {
	b, err := os.ReadFile(t.metadataPath)
	if err != nil && os.IsNotExist(err) {
		err := t.writeLiveVersion(&LiveVersion{Version: -1})
		if err != nil {
			return nil, err
		}
		return &LiveVersion{Version: -1}, nil
	} else if err != nil {
		return nil, err
	}

	var v LiveVersion
	err = yaml.Unmarshal(b, &v)
	return &v, err
}

func (t *MigrationTool) writeLiveVersion(v *LiveVersion) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(t.metadataPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
