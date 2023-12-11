package main

import (
	"flag"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"

	"github.com/vertex-center/vertex/apps"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/netcap"
	"github.com/vertex-center/vlog"
)

func main() {
	defer log.Default.Close()

	ensureRoot()
	parseArgs()

	// If go.mod is there, build vertex first.
	_, err := os.Stat("go.mod")
	if err == nil {
		log.Info("init.go found. Building vertex...")
		buildVertex()
	}

	err = netcap.AllowPortsManagement("vertex")
	if err != nil {
		log.Error(err)
	}

	app.RunKernelApps(apps.Apps)

	exitVertexChan := make(chan error)

	var vertex *exec.Cmd
	go func() {
		defer close(exitVertexChan)

		var err error
		vertex, err = runVertex()
		if err != nil {
			exitVertexChan <- err
			return
		}
		exitVertexChan <- vertex.Wait()
	}()

	for err = range exitVertexChan {
		if err != nil {
			log.Error(err)
		}
	}
}

func ensureRoot() {
	if os.Getuid() != 0 {
		log.Warn("vertex-kernel must be run as root to work properly")
	}
}

func parseArgs() {
	var (
		flagUsername = flag.String("user", "", "username of the unprivileged user")
		flagUID      = flag.Uint("uid", 0, "uid of the unprivileged user")
		flagGID      = flag.Uint("gid", 0, "gid of the unprivileged user")
	)

	flag.Parse()

	if *flagUsername == "" {
		*flagUsername = os.Getenv("USER")
		log.Warn("no username specified; trying to retrieve username from env", vlog.String("user", *flagUsername))
	}

	if *flagUsername != "" {
		u, err := user.Lookup(*flagUsername)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		uid, err := strconv.ParseInt(u.Uid, 10, 32)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		gid, err := strconv.ParseInt(u.Gid, 10, 32)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		config.KernelCurrent.Uid = uint32(uid)
		config.KernelCurrent.Gid = uint32(gid)
		return
	}

	config.KernelCurrent.Uid = uint32(*flagUID)
	config.KernelCurrent.Gid = uint32(*flagGID)
}

func buildVertex() {
	log.Info("Building vertex")

	start := time.Now()
	cmd := exec.Command("go", "build", "-o", "vertex", "cmd/main/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	end := time.Now()

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Build completed in " + end.Sub(start).String())
}
