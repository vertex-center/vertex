package log

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var defaultAgent *Agent

func init() {
	var err error
	defaultAgent, err = NewAgent()
	if err != nil {
		panic(err)
	}

	go func() {
		u := &url.URL{
			Scheme: "ws",
			Host:   "localhost:7516",
			Path:   "/api/logs/ws",
		}
		if err != nil {
			panic(err)
		}

		for {
			err := defaultAgent.Start(u)
			if err != nil {
				err = fmt.Errorf("start log agent: %w", err)
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
				_, _ = fmt.Fprintln(os.Stderr, "retrying in 5 seconds...")
				<-time.After(5 * time.Second)
			}
		}
	}()
}

// Agent is a log agent that gathers logs and buffers them before sending them to Vertex Logs when available.
type Agent struct {
	r, w *os.File
	mu   sync.Mutex // protects w
}

func NewAgent() (*Agent, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	return &Agent{
		r: r,
		w: w,
	}, nil
}

// Start connects to Vertex Logs and starts sending logs.
// Make sure the URL is a websocket URL (ws:// or wss://).
func (a *Agent) Start(u *url.URL) error {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	scanner := bufio.NewScanner(a.r)
	for scanner.Scan() {
		err = conn.WriteMessage(websocket.TextMessage, scanner.Bytes())
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (a *Agent) Close() error {
	var errs []error
	errs = append(errs, a.r.Close())
	errs = append(errs, a.w.Close())
	return errors.Join(errs...)
}

func (a *Agent) Send(s string) (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err = a.w.WriteString(s + "\n")
	return err
}
