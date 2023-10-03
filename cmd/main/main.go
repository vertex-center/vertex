package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/migration"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
)

// version, commit and date will be overridden by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	defer log.Default.Close()

	log.Info("Vertex starting...")

	err := migration.NewMigrationTool(storage.Path).Migrate()
	if err != nil {
		panic(err)
	}

	parseArgs()

	checkNotRoot()

	err = setupDependencies()
	if err != nil {
		log.Error(fmt.Errorf("failed to setup dependencies: %v", err))
		return
	}

	err = config.Current.Apply()
	if err != nil {
		log.Error(fmt.Errorf("failed to apply the current configuration: %v", err))
		return
	}

	r := router.NewRouter(types.About{
		Version: version,
		Commit:  commit,
		Date:    date,

		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	})
	defer r.Stop()

	// Vertex-Kernel Proxy
	var rProxy router.ProxyRouter
	go func() {
		rProxy = router.NewProxyRouter()
		err := rProxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	// Logs
	url := fmt.Sprintf("http://%s", config.Current.HostVertex)
	fmt.Printf("\n-- Vertex Client :: %s\n\n", url)

	r.Start(fmt.Sprintf(":%s", config.Current.PortVertex))
}

func parseArgs() {
	flagVersion := flag.Bool("version", false, "Print vertex version")
	flagV := flag.Bool("v", false, "Print vertex version")
	flagDate := flag.Bool("date", false, "Print the release date")
	flagCommit := flag.Bool("commit", false, "Print the commit hash")

	flagPort := flag.String("port", config.Current.PortVertex, "The Vertex port")
	flagHost := flag.String("host", config.Current.HostVertex, "The Vertex access url")

	flagHostKernel := flag.String("host-kernel", config.Current.HostKernel, "The Vertex Kernel access url")
	flagPortKernel := flag.String("port-kernel", config.Current.PortKernel, "The Vertex Kernel port")

	flagHostProxy := flag.String("host-proxy", config.Current.HostProxy, "The Vertex Proxy access url")
	flagPortProxy := flag.String("port-proxy", config.Current.PortProxy, "The Vertex Proxy port")

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
	config.Current.HostVertex = *flagHost
	config.Current.PortVertex = *flagPort
	config.Current.HostKernel = *flagHostKernel
	config.Current.PortKernel = *flagPortKernel
	config.Current.HostProxy = *flagHostProxy
	config.Current.PortProxy = *flagPortProxy
}

func checkNotRoot() {
	if os.Getuid() == 0 {
		log.Warn("while vertex-kernel must be run as root, the vertex user should not be root")
	}
}

func setupDependencies() error {
	for _, dep := range services.Dependencies {
		err := setupDependency(dep.Updater)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
	return nil
}

func setupDependency(dep types.DependencyUpdater) error {
	err := os.Mkdir(dep.GetPath(), os.ModePerm)
	if os.IsExist(err) {
		// The dependency already exists.
		return nil
	}
	if err != nil {
		return err
	}

	// download
	_, err = dep.CheckForUpdate(types.SettingsUpdatesChannelStable)
	if err != nil && !errors.Is(err, services.ErrDependencyNotInstalled) {
		return err
	}
	return dep.InstallUpdate()
}