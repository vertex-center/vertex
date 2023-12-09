package main

import (
	"github.com/vertex-center/vertex/apps/auth"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandalone(auth.NewApp())
}
