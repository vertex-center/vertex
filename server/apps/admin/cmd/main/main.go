package main

import (
	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/config"
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
	app.RunStandalone(admin.NewApp(), about, true)
}
