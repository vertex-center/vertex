//go:build windows

package reboot

func Reboot() error {
	cmd := exec.Command("shutdown", "/r", "/t", "0")
	return cmd.Run()
}
