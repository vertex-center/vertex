//go:build !windows

package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func allowPort80() {
	cmd := exec.Command("setcap", "cap_net_bind_service=+ep", "vertex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error(errors.New("error trying to allow ./vertex to use the port 80"),
			vlog.String("msg", err.Error()),
		)
	}
}
