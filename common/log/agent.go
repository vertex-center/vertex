package log

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var agent *Agent

func SetupAgent(url url.URL) {
	if strings.HasSuffix(os.Args[0], ".test") {
		return
	}

	var err error
	agent, err = NewAgent()
	if err != nil {
		panic(err)
	}

	go func() {
		url.Scheme = "ws"
		url.Path = "/api/logs/ws"

		_, _ = fmt.Fprintln(os.Stderr, "starting log agent...")
		_, _ = fmt.Fprintln(os.Stderr, "connecting to", url.String())

		for {
			err := agent.Start(&url)
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
// Make sure the Addr is a websocket Addr (ws:// or wss://).
func (a *Agent) Start(u *url.URL) error {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial to %s: %w", u.String(), err)
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
