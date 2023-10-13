package adapter

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	containerstypes "github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
)

type ContainerFilePath string

var (
	ErrContainerNotFound = errors.New("container not found")
)

type ContainerFSAdapter struct {
	containersPath string
}

type ContainerFSAdapterParams struct {
	containersPath string
}

func NewContainerFSAdapter(params *ContainerFSAdapterParams) containerstypes.ContainerAdapterPort {
	if params == nil {
		params = &ContainerFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.Path, "apps", "vx-containers")
	}

	adapter := &ContainerFSAdapter{
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

func (a *ContainerFSAdapter) GetPath(uuid uuid.UUID) string {
	return path.Join(a.containersPath, uuid.String())
}

func (a *ContainerFSAdapter) Create(uuid uuid.UUID) error {
	err := os.MkdirAll(a.GetPath(uuid), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}
	return nil
}

func (a *ContainerFSAdapter) Delete(uuid uuid.UUID) error {
	err := os.RemoveAll(a.GetPath(uuid))
	if err != nil {
		return fmt.Errorf("failed to delete server: %v", err)
	}
	return nil
}

func (a *ContainerFSAdapter) GetAll() ([]uuid.UUID, error) {
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
