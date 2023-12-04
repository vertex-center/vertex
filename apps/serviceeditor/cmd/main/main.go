package main

import (
	"github.com/vertex-center/vertex/apps/serviceeditor"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(serviceeditor.NewApp())
}
