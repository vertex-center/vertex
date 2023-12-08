package main

import (
	"flag"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"

	"github.com/vertex-center/vertex/apps"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/netcap"
	"github.com/vertex-center/vlog"
	"github.com/wI2L/fizz/openapi"
)

var (
	srv *server.Server
	ctx *types.VertexContext
)

func main() {
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

	ctx = types.NewVertexContext(types.About{}, true)
	url := config.Current.KernelURL("vertex")

	info := openapi.Info{
		Title:       "Vertex Kernel",
		Description: "Create your self-hosted lab in one click.",
		Version:     ctx.About().Version,
	}

	srv = server.New("kernel", &info, url, ctx)
	initServices()

	ctx.DispatchEvent(types.EventServerLoad{})
	ctx.DispatchEvent(types.EventServerStart{})

	exitKernelChan := srv.StartAsync()
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

	for {
		select {
		case err := <-exitKernelChan:
			if err != nil {
				log.Error(err)
			}
			if vertex != nil && vertex.Process != nil {
				_ = vertex.Process.Signal(os.Interrupt)
				_, _ = vertex.Process.Wait()
			}
		case err := <-exitVertexChan:
			if err != nil {
				log.Error(err)
			}
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

func initServices() {
	service.NewAppsService(ctx, true, apps.Apps)
}
