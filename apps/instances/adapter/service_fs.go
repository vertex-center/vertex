package adapter

import (
	"errors"
	"fmt"
	"os"
	"path"

	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"gopkg.in/yaml.v3"
)

type ServiceFSAdapter struct {
	servicesPath string
	services     []instancestypes.Service
}

type ServiceFSAdapterParams struct {
	servicesPath string
}

func NewServiceFSAdapter(params *ServiceFSAdapterParams) instancestypes.ServiceAdapterPort {
	if params == nil {
		params = &ServiceFSAdapterParams{}
	}
	if params.servicesPath == "" {
		params.servicesPath = path.Join(storage.Path, "services")
	}

	adapter := &ServiceFSAdapter{
		servicesPath: params.servicesPath,
	}
	err := adapter.Reload()
	if err != nil {
		log.Error(fmt.Errorf("failed to reload services: %v", err))
	}
	return adapter
}

func (a *ServiceFSAdapter) Get(id string) (instancestypes.Service, error) {
	for _, service := range a.services {
		if service.ID == id {
			return service, nil
		}
	}

	return instancestypes.Service{}, instancestypes.ErrServiceNotFound
}

func (a *ServiceFSAdapter) GetRaw(id string) (interface{}, error) {
	servicePath := path.Join(a.servicesPath, "services", id, "service.yml")

	data, err := os.ReadFile(servicePath)
	if err != nil && os.IsNotExist(err) {
		return nil, instancestypes.ErrServiceNotFound
	} else if err != nil {
		return nil, err
	}

	var service interface{}
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *ServiceFSAdapter) GetScript(id string) ([]byte, error) {
	service, err := a.Get(id)
	if err != nil {
		return nil, err
	}

	if service.Methods.Script == nil {
		return nil, errors.New("the service doesn't have a script method")
	}

	return os.ReadFile(path.Join(a.servicesPath, "services", id, service.Methods.Script.Filename))
}

func (a *ServiceFSAdapter) GetAll() []instancestypes.Service {
	return a.services
}

func (a *ServiceFSAdapter) Reload() error {
	servicesPath := path.Join(a.servicesPath, "services")

	a.services = []instancestypes.Service{}

	entries, err := os.ReadDir(servicesPath)
	if err != nil {
		return err
	}

	for _, dir := range entries {
		if !dir.IsDir() {
			continue
		}

		servicePath := path.Join(servicesPath, dir.Name(), "service.yml")

		file, err := os.ReadFile(servicePath)
		if err != nil {
			return err
		}

		var service instancestypes.Service
		err = yaml.Unmarshal(file, &service)
		if err != nil {
			return err
		}

		a.services = append(a.services, service)
	}

	return nil
}
