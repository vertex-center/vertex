package main

import (
	"log"
	"os"

	"github.com/quentinguidee/installer-service/routers"
)

func main() {
	r := routers.InitializeRouter()

	err := os.Mkdir("servers", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Couldn't create 'servers' directory: %v", err)
	}

	err = r.Run(":6130")
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}
