package main

import (
	"github.com/vertex-center/vertex/apps/auth"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(auth.NewApp())
}
