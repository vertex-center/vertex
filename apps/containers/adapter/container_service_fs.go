package adapter

import (
	"os"
	"path"

	"github.com/google/uuid"
	containerstypes "github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

const (
	ContainerServicePath = ".vertex/service.yml"
)

type ContainerServiceFSAdapter struct {
	containersPath string
}

type ContainerServiceFSAdapterParams struct {
	containersPath string
}

func NewContainerServiceFSAdapter(params *ContainerServiceFSAdapterParams) containerstypes.ContainerServiceAdapterPort {
	if params == nil {
		params = &ContainerServiceFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.Path, "apps", "vx-containers")
	}

	adapter := &ContainerServiceFSAdapter{
		containersPath: params.containersPath,
	}

	return adapter
}

func (a *ContainerServiceFSAdapter) Save(uuid uuid.UUID, service containerstypes.Service) error {
	servicePath := path.Join(a.containersPath, uuid.String(), ContainerServicePath)

	err := os.MkdirAll(path.Join(a.containersPath, uuid.String(), ".vertex"), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	serviceBytes, err := yaml.Marshal(service)
	if err != nil {
		return err
	}

	err = os.WriteFile(servicePath, serviceBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (a *ContainerServiceFSAdapter) Load(uuid uuid.UUID) (containerstypes.Service, error) {
	servicePath := path.Join(a.containersPath, uuid.String(), ContainerServicePath)

	data, err := os.ReadFile(servicePath)
	if err != nil {
		log.Warn("file not found", vlog.String("path", servicePath))
		return containerstypes.Service{}, err
	}

	var service containerstypes.Service
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *ContainerServiceFSAdapter) LoadRaw(uuid uuid.UUID) (interface{}, error) {
	servicePath := path.Join(a.containersPath, uuid.String(), ContainerServicePath)

	data, err := os.ReadFile(servicePath)
	if err != nil {
		return nil, err
	}

	var service interface{}
	err = yaml.Unmarshal(data, &service)
	return service, err
}
