package main

import (
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(tunnels.NewApp())
}
