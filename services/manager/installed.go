package servicesmanager

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services/instance"
	"github.com/vertex-center/vertex/storage"
)

var logger = console.New("vertex::services-manager")

func ReloadAllInstalled() error {
	entries, err := os.ReadDir(storage.PathInstances)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		isInstance := entry.IsDir() || info.Mode()&os.ModeSymlink != 0

		if isInstance {
			logger.Log(fmt.Sprintf("found service uuid=%s", entry.Name()))
			serviceUUID, err := uuid.Parse(entry.Name())
			if err != nil {
				return err
			}

			if !instance.Exists(serviceUUID) {
				logger.Log(fmt.Sprintf("instantiate service uuid=%s", entry.Name()))

				_, err = instance.Instantiate(serviceUUID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
