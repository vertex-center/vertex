//go:build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func runVertex() (*exec.Cmd, error) {
	uid, gid := config.KernelCurrent.Uid, config.KernelCurrent.Gid

	log.Info("running vertex",
		vlog.Uint32("uid", uid),
		vlog.Uint32("gid", gid),
	)

	cmd := exec.Command("./vertex", "-port", config.Current.Port, "-host", config.Current.Host)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
}
