package main

import (
	"fmt"
	"os"

	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/router"
)

var logger = console.New("vertex")

func main() {
	r := router.InitializeRouter()

	err := os.Mkdir("servers", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Error(fmt.Errorf("couldn't create 'servers' directory: %v", err))
		return
	}

	err = os.Mkdir("clients", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Error(fmt.Errorf("couldn't create 'clients' directory: %v", err))
		return
	}

	err = r.Run(":6130")
	if err != nil {
		logger.Error(fmt.Errorf("error while starting server: %v", err))
		return
	}
}
