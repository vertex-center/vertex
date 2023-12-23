package main

import (
	"os"

	"github.com/vertex-center/vertex/apps"
	logsmeta "github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/config"
)

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	defer log.Default.Close()

	ensureNotRoot()

	about := common.NewAbout(version, commit, date)
	for _, a := range apps.Apps {
		meta := a.Meta()
		config.RegisterHost(meta.ID, meta.DefaultPort)
		config.RegisterHost(meta.ID+"-kernel", meta.DefaultKernelPort)
	}
	config.RegisterHost("vertex", "6130")
	config.ParseArgs(about)

	log.SetupAgent(*config.Current.Addr(logsmeta.Meta.ID))

	app.RunApps(about, apps.Apps)
}

func ensureNotRoot() {
	if os.Getuid() == 0 {
		log.Warn("while vertex-kernel must be run as root, the vertex user should not be root")
	}
}
