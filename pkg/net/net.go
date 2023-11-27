package net

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
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
	ch := make(chan bool)
	timeout := time.After(10 * time.Second)

	go func() {
		for {
			_, err := http.Get(url)
			if err == nil {
				ch <- true
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {
	case <-ch:
		log.Info("successfully pinged", vlog.String("url", url))
		return nil
	case <-timeout:
		return fmt.Errorf("internet connection: Failed to ping %s", url)
	}
}

func WaitInternetConn() error {
	return Wait("https://www.google.com/robots.txt")
}
