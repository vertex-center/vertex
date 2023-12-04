package main

import (
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(reverseproxy.NewApp())
}
