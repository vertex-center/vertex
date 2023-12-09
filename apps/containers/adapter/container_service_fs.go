package adapter

import (
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/cmd/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

const (
	ContainerServicePath = ".vertex/service.yml"
)

type containerServiceFSAdapter struct {
	containersPath string
}

type ContainerServiceFSAdapterParams struct {
	containersPath string
}

func NewContainerServiceFSAdapter(params *ContainerServiceFSAdapterParams) port.ContainerServiceAdapter {
	if params == nil {
		params = &ContainerServiceFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	adapter := &containerServiceFSAdapter{
		containersPath: params.containersPath,
	}

	return adapter
}

func (a *containerServiceFSAdapter) Save(uuid uuid.UUID, service types.Service) error {
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

func (a *containerServiceFSAdapter) Load(uuid uuid.UUID) (types.Service, error) {
	servicePath := path.Join(a.containersPath, uuid.String(), ContainerServicePath)

	data, err := os.ReadFile(servicePath)
	if err != nil {
		log.Warn("file not found", vlog.String("path", servicePath))
		return types.Service{}, err
	}

	var service types.Service
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *containerServiceFSAdapter) LoadRaw(uuid uuid.UUID) (interface{}, error) {
	servicePath := path.Join(a.containersPath, uuid.String(), ContainerServicePath)

	data, err := os.ReadFile(servicePath)
	if err != nil {
		return nil, err
	}

	var service interface{}
	err = yaml.Unmarshal(data, &service)
	return service, err
}
