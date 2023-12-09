package main

import (
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandalone(tunnels.NewApp())
}
