package main

import (
	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandaloneKernel(admin.NewApp())
}
