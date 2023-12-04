package main

import (
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(containers.NewApp())
}
