package net

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
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

func Wait(ctx context.Context, url string) error {
	ch := make(chan bool)
	done := ctx.Done()

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
		return nil
	case <-done:
		return fmt.Errorf("failed to ping %s", url)
	}
}

func WaitInternetConn(ctx context.Context) error {
	return Wait(ctx, "https://clients3.google.com/generate_204")
}

func WaitInternetConnWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return WaitInternetConn(ctx)
}
