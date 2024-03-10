package main

import (
	"github.com/vertex-center/vertex/server/apps/tunnels"
	"github.com/vertex-center/vertex/server/common"
	"github.com/vertex-center/vertex/server/common/app"
)

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	about := common.NewAbout(version, commit, date)
	app.RunStandalone(tunnels.NewApp(), about, true)
}
