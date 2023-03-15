package servicesmanager

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services/instances"
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

			if !instances.Exists(serviceUUID) {
				logger.Log(fmt.Sprintf("instantiate service uuid=%s", entry.Name()))

				_, err = instances.Instantiate(serviceUUID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
