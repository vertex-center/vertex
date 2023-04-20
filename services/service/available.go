package service

import (
	"os"
	"path"

	"github.com/goccy/go-json"
	"github.com/vertex-center/vertex/storage"
)

var available []Service

func ReloadAvailableServices() error {
	return reloadAvailableServices(storage.PathServices)
}

func reloadAvailableServices(servicesPath string) error {
	url := "https://github.com/vertex-center/vertex-services"

	err := storage.CloneOrPullRepository(url, servicesPath)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(path.Join(servicesPath, "services.json"))
	if err != nil {
		return err
	}

	var availableMap map[string]Service
	err = json.Unmarshal(file, &availableMap)
	if err != nil {
		return err
	}

	available = []Service{}
	for key, service := range availableMap {
		service.ID = key
		available = append(available, service)
	}

	return nil
}

func ListAvailable() []Service {
	return available
}
