package main

import (
	"errors"
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
	"github.com/vertex-center/vlog"
)

func main() {
	ensureRoot()

	parseArgs()

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
		vertex, err = runVertex()
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
	_ = r.Stop()
	if vertex != nil && vertex.Process != nil {
		_ = vertex.Process.Kill()
	}
}

func ensureRoot() {
	if os.Getuid() != 0 {
		log.Warn("vertex-kernel must be run as root to work properly")
	}
}

func parseArgs() {
	flagUsername := flag.String("user", "", "username of the unprivileged user")
	flagUID := flag.Uint("uid", 0, "uid of the unprivileged user")
	flagGID := flag.Uint("gid", 0, "gid of the unprivileged user")

	// Copy/paste from cmd/main/main.go
	flagPort := flag.String("port", config.Current.Port, "The Vertex port")
	flagHost := flag.String("host", config.Current.Host, "The Vertex access url")

	flag.Parse()

	config.Current.Host = *flagHost
	config.Current.Port = *flagPort

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

func runVertex() (*exec.Cmd, error) {
	uid, gid := config.KernelCurrent.Uid, config.KernelCurrent.Gid

	// If go.mod is there, build vertex first.
	_, err := os.Stat("go.mod")
	if err == nil {
		log.Info("init.go found. Building vertex...")
		buildVertex()
	}

	// Allow Vertex Proxy to use the port 80
	cmd := exec.Command("setcap", "cap_net_bind_service=+ep", "vertex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Error(errors.New("error trying to allow ./vertex to use the port 80"),
			vlog.String("msg", err.Error()),
		)
	}

	// Run Vertex
	log.Info("running vertex",
		vlog.Uint32("uid", uid),
		vlog.Uint32("gid", gid),
	)

	cmd = exec.Command("./vertex", "-port", config.Current.Port, "-host", config.Current.Host)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
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
