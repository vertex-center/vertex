package main

import (
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandalone(reverseproxy.NewApp())
}
