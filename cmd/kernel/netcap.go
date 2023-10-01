//go:build !windows

package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func allowPortsManagement() {
	cmd := exec.Command("setcap", "cap_net_bind_service=+ep", "vertex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error(errors.New("error trying to allow ./vertex to manage ports"),
			vlog.String("msg", err.Error()),
		)
	}
}
