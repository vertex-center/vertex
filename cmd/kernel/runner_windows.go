//go:build windows

package main

import (
	"os"
	"os/exec"

	"github.com/vertex-center/vertex/config"
)

func runVertex() (*exec.Cmd, error) {
	cmd := exec.Command("vertex.exe", "-port", config.Current.Port, "-host", config.Current.Host)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
}
