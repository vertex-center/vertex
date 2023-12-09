package adapter

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/cmd/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type ContainerFilePath string

var (
	ErrContainerNotFound = errors.New("container not found")
)

type containerFSAdapter struct {
	containersPath string
}

type ContainerFSAdapterParams struct {
	containersPath string
}

func NewContainerFSAdapter(params *ContainerFSAdapterParams) port.ContainerAdapter {
	if params == nil {
		params = &ContainerFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	adapter := &containerFSAdapter{
		containersPath: params.containersPath,
	}

	err := os.MkdirAll(adapter.containersPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", adapter.containersPath),
		)
		os.Exit(1)
	}

	return adapter
}

func (a *containerFSAdapter) GetPath(uuid uuid.UUID) string {
	return path.Join(a.containersPath, uuid.String())
}

func (a *containerFSAdapter) Create(uuid uuid.UUID) error {
	err := os.MkdirAll(a.GetPath(uuid), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	return nil
}

func (a *containerFSAdapter) Delete(uuid uuid.UUID) error {
	err := os.RemoveAll(a.GetPath(uuid))
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}
	return nil
}

func (a *containerFSAdapter) GetAll() ([]uuid.UUID, error) {
	entries, err := os.ReadDir(a.containersPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var uuids []uuid.UUID
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Error(err)
			continue
		}

		isContainer := entry.IsDir() || info.Mode()&os.ModeSymlink != 0

		if isContainer {
			log.Info("found container",
				vlog.String("uuid", entry.Name()),
			)

			id, err := uuid.Parse(entry.Name())
			if err != nil {
				log.Error(err)
				continue
			}

			uuids = append(uuids, id)
		}
	}

	return uuids, nil
}
