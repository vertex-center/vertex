//go:build windows

package netcap

func AllowPortsManagement(execPath string) error {
	// ignored on Windows
	return nil
}
