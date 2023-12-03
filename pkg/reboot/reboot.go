//go:build !windows

package reboot

import "os/exec"

func Reboot() error {
	cmd := exec.Command("reboot")
	return cmd.Run()
}
