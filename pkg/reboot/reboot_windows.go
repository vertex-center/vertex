//go:build windows

package reboot

import "syscall"

func Reboot() error {
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}
