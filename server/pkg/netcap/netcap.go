//go:build !windows

package netcap

import (
	"os"
	"os/exec"
)

func AllowPortsManagement(execPath string) error {
	cmd := exec.Command("setcap", "cap_net_bind_service=+ep", execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
