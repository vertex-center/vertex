package adapter

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"gopkg.in/yaml.v3"
)

type serviceFSAdapter struct {
	servicesPath string
	services     []types.Service
}

type ServiceFSAdapterParams struct {
	servicesPath string
}

func NewServiceFSAdapter(params *ServiceFSAdapterParams) port.ServiceAdapter {
	if params == nil {
		params = &ServiceFSAdapterParams{}
	}
	if params.servicesPath == "" {
		params.servicesPath = path.Join(storage.FSPath, "services")
	}

	adapter := &serviceFSAdapter{
		servicesPath: params.servicesPath,
	}
	err := adapter.Reload()
	if err != nil {
		log.Error(fmt.Errorf("failed to reload services: %w", err))
	}
	return adapter
}

func (a *serviceFSAdapter) Get(id string) (types.Service, error) {
	for _, service := range a.services {
		if service.ID == id {
			return service, nil
		}
	}

	return types.Service{}, types.ErrServiceNotFound
}

func (a *serviceFSAdapter) GetRaw(id string) (interface{}, error) {
	servicePath := path.Join(a.servicesPath, "services", id, "service.yml")

	data, err := os.ReadFile(servicePath)
	if err != nil && os.IsNotExist(err) {
		return nil, types.ErrServiceNotFound
	} else if err != nil {
		return nil, err
	}

	var service interface{}
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *serviceFSAdapter) GetScript(id string) ([]byte, error) {
	service, err := a.Get(id)
	if err != nil {
		return nil, err
	}

	if service.Methods.Script == nil {
		return nil, errors.New("the service doesn't have a script method")
	}

	return os.ReadFile(path.Join(a.servicesPath, "services", id, service.Methods.Script.Filename))
}

func (a *serviceFSAdapter) GetAll() []types.Service {
	return a.services
}

func (a *serviceFSAdapter) Reload() error {
	servicesPath := path.Join(a.servicesPath, "services")

	a.services = []types.Service{}

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

		var service types.Service
		err = yaml.Unmarshal(file, &service)
		if err != nil {
			return err
		}

		a.services = append(a.services, service)
	}

	return nil
}
