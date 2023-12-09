package main

import (
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandalone(monitoring.NewApp())
}
