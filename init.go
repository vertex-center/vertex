package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/client"
	"github.com/vertex-center/vertex/dependency/package"
	"github.com/vertex-center/vertex/router"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/updater"
)

// version, commit and date will be overridden by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var logger = console.New("vertex")

func main() {
	parseArgs()

	err := updater.CheckForUpdates(version)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to check for updates: %v", err))
	}

	err = os.MkdirAll(storage.PathInstances, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Error(fmt.Errorf("couldn't create '%s' directory: %v", storage.PathInstances, err))
		return
	}

	err = client.Setup()
	if err != nil {
		logger.Error(fmt.Errorf("failed to setup the web client: %v", err))
		return
	}

	err = _package.Reload()
	if err != nil {
		logger.Error(fmt.Errorf("failed to reload dependencies: %v", err))
		return
	}

	err = servicesmanager.Reload()
	if err != nil {
		logger.Error(fmt.Errorf("failed to reload available services: %v", err))
		return
	}

	r := router.InitializeRouter(router.About{
		Version: version,
		Commit:  commit,
		Date:    date,
	})

	err = servicesmanager.ReloadAllInstalled()
	if err != nil {
		logger.Error(err)
		return
	}

	err = r.Run(":6130")
	if err != nil {
		logger.Error(fmt.Errorf("error while starting server: %v", err))
		return
	}
}

func parseArgs() {
	flagVersion := flag.Bool("version", false, "Print vertex version")
	flagV := flag.Bool("v", false, "Print vertex version")
	flagDate := flag.Bool("date", false, "Print the release date")
	flagCommit := flag.Bool("commit", false, "Print the commit hash")

	flag.Parse()
	if *flagVersion || *flagV {
		fmt.Println(version)
		os.Exit(0)
	}
	if *flagDate {
		fmt.Println(date)
		os.Exit(0)

	}
	if *flagCommit {
		fmt.Println(commit)
		os.Exit(0)
	}
}
