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
			&migration1{},
		},
	}
}

func (t *MigrationTool) Migrate() ([]interface{}, error) {
	v, err := t.readLiveVersion()
	if err != nil {
		return nil, err
	}

	log.Info("current live directory version", vlog.Int("version", v.Version))

	// migrationCommands is a slice of commands that will be dispatched in
	// Vertex after all migrations are applied.
	var migrationCommands []interface{}

	for i := v.Version + 1; i < len(t.migrations); i++ {
		log.Info("running migration", vlog.Int("version", i))
		err := t.migrations[i].Up(t.livePath)
		if err != nil {
			return nil, err
		}

		// If the migration implements the CommandsDispatcher interface, we
		// append the commands to the migrationCommands slice.
		if c, ok := t.migrations[i].(CommandsDispatcher); ok {
			log.Info("dispatching commands", vlog.Int("version", i))
			migrationCommands = append(migrationCommands, c.DispatchCommands()...)
		}

		v.Version = i
		err = t.writeLiveVersion(v)
		if err != nil {
			return migrationCommands, err
		}
	}

	return migrationCommands, nil
}

func (t *MigrationTool) readLiveVersion() (*LiveVersion, error) {
	b, err := os.ReadFile(t.metadataPath)
	if err != nil && os.IsNotExist(err) {
		liveVersion := &LiveVersion{Version: len(t.migrations) - 1}
		err := t.writeLiveVersion(liveVersion)
		if err != nil {
			return nil, err
		}
		return liveVersion, nil
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
