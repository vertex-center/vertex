package vdocker

import "os"

func RunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}
