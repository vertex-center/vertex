//go:build windows

package main

import (
	"os"
	"os/exec"
)

func runVertex(args ...string) (*exec.Cmd, error) {
	cmd := exec.Command("vertex.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
}
