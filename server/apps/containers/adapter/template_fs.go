package adapter

import (
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/common/storage"
	"gopkg.in/yaml.v3"
)

type templateFSAdapter struct {
	servicesPath string
	templates    []types.Template
}

type TemplateFSAdapterParams struct {
	templatesPath string
}

func NewTemplateFSAdapter(params *TemplateFSAdapterParams) port.TemplateAdapter {
	if params == nil {
		params = &TemplateFSAdapterParams{}
	}
	if params.templatesPath == "" {
		params.templatesPath = path.Join(storage.FSPath, "services")
	}

	adapter := &templateFSAdapter{
		servicesPath: params.templatesPath,
	}
	err := adapter.Reload()
	if err != nil {
		log.Error(fmt.Errorf("failed to reload templates: %w", err))
	}
	return adapter
}

func (a *templateFSAdapter) Get(id string) (types.Template, error) {
	for _, service := range a.templates {
		if service.ID == id {
			return service, nil
		}
	}

	return types.Template{}, types.ErrTemplateNotFound
}

func (a *templateFSAdapter) GetRaw(id string) (interface{}, error) {
	servicePath := path.Join(a.servicesPath, "services", id, "service.yml")

	data, err := os.ReadFile(servicePath)
	if err != nil && os.IsNotExist(err) {
		return nil, types.ErrTemplateNotFound
	} else if err != nil {
		return nil, err
	}

	var template interface{}
	err = yaml.Unmarshal(data, &template)
	return template, err
}

func (a *templateFSAdapter) GetAll() []types.Template {
	return a.templates
}

func (a *templateFSAdapter) Reload() error {
	servicesPath := path.Join(a.servicesPath, "services")

	a.templates = []types.Template{}

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

		var template types.Template
		err = yaml.Unmarshal(file, &template)
		if err != nil {
			return err
		}

		a.templates = append(a.templates, template)
	}

	return nil
}
