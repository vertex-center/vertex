package config

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/vertex-center/vertex/pkg/net"
)

const DefaultApiURLFormat = "http://%s:%s/api"

var Current = New()

type Mode string

const (
	ProductionMode Mode = "production"
	DebugMode      Mode = "debug"
)

type Config struct {
	mode    Mode
	localIP string
	mu      sync.RWMutex

	Port  string
	hosts map[string]string
}

func New() *Config {
	localIP, err := net.LocalIP()
	if err != nil {
		localIP = "127.0.0.1"
	}

	c := &Config{
		mode:    ProductionMode,
		localIP: localIP,
		hosts:   map[string]string{},
	}

	if os.Getenv("DEBUG") == "1" {
		println("debug mode enabled. proceed with caution!")
		c.mode = DebugMode
	}

	return c
}

func (c *Config) Addr(id string) *url.URL {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if u, ok := c.hosts[id]; ok {
		p, err := url.Parse(u)
		if err != nil {
			return &url.URL{}
		}
		return p
	}
	return &url.URL{}
}

func (c *Config) KernelAddr(id string) *url.URL {
	return c.Addr(id + "-kernel")
}

func (c *Config) DefaultApiAddr(defaultPort string) string {
	return fmt.Sprintf(DefaultApiURLFormat, c.localIP, defaultPort)
}

func (c *Config) RegisterAPIAddr(id string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.hosts[id]; ok {
		return
	}
	c.hosts[id] = url
}

func (c *Config) SetAPIAddr(id string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hosts[id] = url
}

func (c *Config) LocalIP() string {
	return c.localIP
}

func (c *Config) Debug() bool {
	return c.mode == DebugMode
}
