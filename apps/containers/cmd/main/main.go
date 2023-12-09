package main

import (
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandalone(containers.NewApp(), true)
}
