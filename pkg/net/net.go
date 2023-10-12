package net

import (
	"fmt"
	"net"
	"time"

	"github.com/antelman107/net-wait-go/wait"
)

func LocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func Wait(url string) error {
	if !wait.New(
		wait.WithWait(time.Second),
		wait.WithBreak(500*time.Millisecond),
	).Do([]string{url}) {
		return fmt.Errorf("internet connection: Failed to ping %s", url)
	} else {
		return nil
	}
}
