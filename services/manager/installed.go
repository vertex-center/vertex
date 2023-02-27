package servicesmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services"
)

var logger = console.New("vertex::services-manager")

func ReloadAllInstalled() error {
	entries, err := os.ReadDir("servers")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			data, err := os.ReadFile(path.Join("servers", entry.Name(), ".vertex", "service.json"))
			if err != nil {
				logger.Warn(fmt.Sprintf("service '%s' has no '.vertex/service.json' file", entry.Name()))
				continue
			}

			var service services.Service
			err = json.Unmarshal(data, &service)
			if err != nil {
				return err
			}

			if !service.IsInstalled() {
				logger.Log(fmt.Sprintf("service found: '%s'", service.ID))
				_, err = service.Instantiate()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
