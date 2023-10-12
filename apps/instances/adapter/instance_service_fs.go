package adapter

import (
	"os"
	"path"

	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

const (
	InstanceServicePath = ".vertex/service.yml"
)

type InstanceServiceFSAdapter struct {
	instancesPath string
}

type InstanceServiceFSAdapterParams struct {
	instancesPath string
}

func NewInstanceServiceFSAdapter(params *InstanceServiceFSAdapterParams) instancestypes.InstanceServiceAdapterPort {
	if params == nil {
		params = &InstanceServiceFSAdapterParams{}
	}
	if params.instancesPath == "" {
		params.instancesPath = path.Join(storage.Path, "instances")
	}

	adapter := &InstanceServiceFSAdapter{
		instancesPath: params.instancesPath,
	}

	return adapter
}

func (a *InstanceServiceFSAdapter) Save(uuid uuid.UUID, service instancestypes.Service) error {
	servicePath := path.Join(a.instancesPath, uuid.String(), InstanceServicePath)

	err := os.MkdirAll(path.Join(a.instancesPath, uuid.String(), ".vertex"), os.ModePerm)
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

func (a *InstanceServiceFSAdapter) Load(uuid uuid.UUID) (instancestypes.Service, error) {
	servicePath := path.Join(a.instancesPath, uuid.String(), InstanceServicePath)

	data, err := os.ReadFile(servicePath)
	if err != nil {
		log.Warn("file not found", vlog.String("path", servicePath))
		return instancestypes.Service{}, err
	}

	var service instancestypes.Service
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *InstanceServiceFSAdapter) LoadRaw(uuid uuid.UUID) (interface{}, error) {
	servicePath := path.Join(a.instancesPath, uuid.String(), InstanceServicePath)

	data, err := os.ReadFile(servicePath)
	if err != nil {
		return nil, err
	}

	var service interface{}
	err = yaml.Unmarshal(data, &service)
	return service, err
}
