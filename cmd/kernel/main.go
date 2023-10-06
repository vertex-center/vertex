package main

import (
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/router"
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

	allowPortsManagement()

	shutdownChan := make(chan os.Signal, 1)

	// Vertex-Kernel
	var r router.KernelRouter
	go func() {
		r = router.NewKernelRouter()
		err := r.Start()
		if err != nil {
			log.Error(err)
		}
		shutdownChan <- syscall.SIGINT
	}()

	// Vertex
	var vertex *exec.Cmd
	go func() {
		var err error
		vertex, err = runVertex([]string{
			"-host", config.KernelCurrent.Host,
			"-port", config.KernelCurrent.Port,
			"-port-kernel", config.KernelCurrent.PortKernel,
			"-port-proxy", config.KernelCurrent.PortProxy,
			"-port-prometheus", config.KernelCurrent.PortPrometheus,
		}...)
		if err != nil {
			log.Error(err)
			shutdownChan <- syscall.SIGINT
			return
		}
		err = vertex.Wait()
		if err != nil {
			log.Error(err)
		}
		shutdownChan <- syscall.SIGINT
	}()

	// OS interrupt
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	<-shutdownChan
	log.Info("Shutting down...")
	if vertex != nil && vertex.Process != nil {
		_ = vertex.Process.Signal(os.Interrupt)
		_, _ = vertex.Process.Wait()
	}
	_ = r.Stop()
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

		flagHost = flag.String("host", config.Current.Host, "The Vertex access url")

		flagPort           = flag.String("port", config.Current.Port, "The Vertex port")
		flagPortKernel     = flag.String("port-kernel", config.Current.PortKernel, "The Vertex Kernel port")
		flagPortProxy      = flag.String("port-proxy", config.Current.PortProxy, "The Vertex Proxy port")
		flagPortPrometheus = flag.String("port-prometheus", config.Current.PortPrometheus, "The Prometheus port")
	)

	flag.Parse()

	config.KernelCurrent.Host = *flagHost
	config.KernelCurrent.Port = *flagPort
	config.KernelCurrent.PortKernel = *flagPortKernel
	config.KernelCurrent.PortProxy = *flagPortProxy
	config.KernelCurrent.PortPrometheus = *flagPortPrometheus

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
	cmd := exec.Command("go", "build", "-o", "vertex", "cmd/main/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
