package main

import (
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(monitoring.NewApp())
}
