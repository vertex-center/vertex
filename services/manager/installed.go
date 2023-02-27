package servicesmanager

import (
	"fmt"
	"os"

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
			service, err := services.FromDisk(entry.Name())
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
