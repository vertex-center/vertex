package main

import (
	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandaloneKernel(admin.NewApp())
}
