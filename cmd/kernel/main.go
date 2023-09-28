package main

import (
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/router"
	"github.com/vertex-center/vlog"
)

func main() {
	ensureRoot()

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
		u := getUnprivilegedUser()
		var err error
		vertex, err = runVertex(u)
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

func getUnprivilegedUser() *user.User {
	flagUnprivilegedUsername := flag.String("u", "vertex", "unprivileged username")
	flag.Parse()
	username := *flagUnprivilegedUsername
	return getUser(username)
}

func getUser(username string) *user.User {
	u, err := user.Lookup(username)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return u
}

func runVertex(user *user.User) (*exec.Cmd, error) {
	uid, err := strconv.ParseInt(user.Uid, 10, 32)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.ParseInt(user.Gid, 10, 32)
	if err != nil {
		return nil, err
	}

	// If vertex init.go is there, build vertex first.
	_, err = os.Stat("init.go")
	if err == nil {
		log.Info("init.go found. Building vertex...")
		buildVertex()
	}

	// Run Vertex
	log.Info("running vertex",
		vlog.String("as", user.Username),
		vlog.Int64("uid", uid),
		vlog.Int64("gid", gid),
	)

	cmd := exec.Command("./vertex")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
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
