package main

import (
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app"
)

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	about := common.NewAbout(version, commit, date)
	app.RunStandalone(reverseproxy.NewApp(), about, true)
}
