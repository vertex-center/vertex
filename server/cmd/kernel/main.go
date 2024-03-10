package main

import (
	"os"

	"github.com/vertex-center/vertex/server/apps"
	logsmeta "github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/common"
	"github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/config"
)

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	defer log.Default.Close()
	ensureRoot()

	about := common.NewAbout(version, commit, date)
	for _, a := range apps.Apps {
		meta := a.Meta()
		config.RegisterHost(meta.ID, meta.DefaultPort)
		config.RegisterHost(meta.ID+"-kernel", meta.DefaultKernelPort)
	}
	config.ParseArgs(about)

	log.SetupAgent(*config.Current.Addr(logsmeta.Meta.ID))

	app.RunKernelApps(about, apps.Apps)
}

func ensureRoot() {
	if os.Getuid() != 0 {
		log.Warn("vertex-kernel must be run as root to work properly")
	}
}
