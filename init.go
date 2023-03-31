package main

import (
	"fmt"
	"os"

	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/client"
	"github.com/vertex-center/vertex/dependency/package"
	"github.com/vertex-center/vertex/router"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
	"github.com/vertex-center/vertex/storage"
)

var logger = console.New("vertex")

func main() {
	err := os.MkdirAll(storage.PathInstances, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Error(fmt.Errorf("couldn't create '%s' directory: %v", storage.PathInstances, err))
		return
	}

	err = client.Setup()
	if err != nil {
		logger.Error(fmt.Errorf("failed to setup the web client: %v", err))
		return
	}

	err = _package.Reload()
	if err != nil {
		logger.Error(fmt.Errorf("failed to reload dependencies: %v", err))
		return
	}

	r := router.InitializeRouter()

	err = servicesmanager.ReloadAllInstalled()
	if err != nil {
		logger.Error(err)
		return
	}

	err = r.Run(":6130")
	if err != nil {
		logger.Error(fmt.Errorf("error while starting server: %v", err))
		return
	}
}
