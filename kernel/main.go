package main

import (
	"flag"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/router"
	"github.com/vertex-center/vlog"
)

func main() {
	ensureRoot()

	r := router.NewKernelRouter()

	go func() {
		err := r.Start()
		defer r.Stop()

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()

	u := getUnprivilegedUser()

	go func() {
		vertex, err := runVertex(u)
		defer vertex.Process.Kill()

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = vertex.Wait()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()

	select {}
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
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return nil, err
	}

	log.Info("Running vertex",
		vlog.String("username", user.Username),
		vlog.Int("uid", uid),
		vlog.Int("gid", gid),
	)

	cmd := exec.Command("./vertex")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
}