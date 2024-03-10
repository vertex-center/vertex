package main

import (
	"github.com/vertex-center/vertex/server/apps/auth"
	"github.com/vertex-center/vertex/server/common"
	"github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/config"
)

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	about := common.NewAbout(version, commit, date)
	config.Current.RegisterDBArgs()
	app.RunStandalone(auth.NewApp(), about, true)
}
