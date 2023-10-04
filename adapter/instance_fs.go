package adapter

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type InstanceFilePath string

var (
	ErrContainerNotFound = errors.New("container not found")
)

type InstanceFSAdapter struct {
	instancesPath string
}

type InstanceFSAdapterParams struct {
	instancesPath string
}

func NewInstanceFSAdapter(params *InstanceFSAdapterParams) types.InstanceAdapterPort {
	if params == nil {
		params = &InstanceFSAdapterParams{}
	}
	if params.instancesPath == "" {
		params.instancesPath = path.Join(storage.Path, "instances")
	}

	adapter := &InstanceFSAdapter{
		instancesPath: params.instancesPath,
	}

	err := os.MkdirAll(adapter.instancesPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", adapter.instancesPath),
		)
		os.Exit(1)
	}

	return adapter
}

func (a *InstanceFSAdapter) GetPath(uuid uuid.UUID) string {
	return path.Join(a.instancesPath, uuid.String())
}

func (a *InstanceFSAdapter) Delete(uuid uuid.UUID) error {
	err := os.RemoveAll(a.GetPath(uuid))
	if err != nil {
		return fmt.Errorf("failed to delete server: %v", err)
	}
	return nil
}

func (a *InstanceFSAdapter) GetAll() ([]uuid.UUID, error) {
	entries, err := os.ReadDir(a.instancesPath)
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

		isInstance := entry.IsDir() || info.Mode()&os.ModeSymlink != 0

		if isInstance {
			log.Info("found instance",
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
