package servicesmanager

import (
	"fmt"
	"os"

	"github.com/google/uuid"
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
			logger.Log(fmt.Sprintf("found service uuid=%s", entry.Name()))
			serviceUUID, err := uuid.Parse(entry.Name())
			if err != nil {
				return err
			}

			if !services.IsInstantiated(serviceUUID) {
				logger.Log(fmt.Sprintf("instantiate service uuid=%s", entry.Name()))

				_, err = services.Instantiate(serviceUUID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
