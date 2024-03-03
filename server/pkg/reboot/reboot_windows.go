//go:build windows

package reboot

import "os/exec"

func Reboot() error {
	cmd := exec.Command("shutdown", "/r", "/t", "0")
	return cmd.Run()
}
