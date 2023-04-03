package servicesmanager

import (
	"os"
	"path"

	"github.com/goccy/go-json"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/storage"
)

var available []services.Service

func Reload() error {
	return reload(storage.PathServices)
}

func reload(servicesPath string) error {
	url := "https://github.com/vertex-center/vertex-services"

	err := storage.CloneOrPullRepository(url, servicesPath)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(path.Join(servicesPath, "services.json"))
	if err != nil {
		return err
	}

	var availableMap map[string]services.Service
	err = json.Unmarshal(file, &availableMap)
	if err != nil {
		return err
	}

	available = []services.Service{}
	for key, service := range availableMap {
		service.ID = key
		available = append(available, service)
	}

	return nil
}

func ListAvailable() []services.Service {
	return available
}
