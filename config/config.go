package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
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
	mode       Mode
	localIP    string
	kernelUrls map[string]string
	mu         sync.RWMutex
	Urls       map[string]string
}

func New() *Config {
	localIP, err := net.LocalIP()
	if err != nil {
		localIP = "127.0.0.1"
	}

	c := &Config{
		mode:    ProductionMode,
		localIP: localIP,
		Urls: map[string]string{
			"vertex": fmt.Sprintf(DefaultApiURLFormat, localIP, "6130"),
		},
		kernelUrls: map[string]string{
			"vertex": fmt.Sprintf(DefaultApiURLFormat, localIP, "6131"),
		},
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	env := os.Environ()
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			key, value := pair[0], pair[1]
			if strings.HasPrefix(key, "VERTEX_URL_") {
				name := strings.TrimPrefix(key, "VERTEX_URL_")
				name = strings.ToLower(name)
				if strings.HasSuffix(key, "_kernel") {
					name = strings.TrimSuffix(name, "_kernel")
					c.kernelUrls[name] = value
				} else {
					c.Urls[name] = value
				}
			}
		}
	}

	if os.Getenv("DEBUG") == "1" {
		println("debug mode enabled. proceed with caution!")
		c.mode = DebugMode
	}

	return c
}

func (c *Config) KernelURL(id string) *url.URL {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if u, ok := c.kernelUrls[id]; ok {
		p, err := url.Parse(u)
		if err != nil {
			return &url.URL{}
		}
		return p
	}
	return &url.URL{}
}

func (c *Config) URL(id string) *url.URL {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if u, ok := c.Urls[id]; ok {
		p, err := url.Parse(u)
		if err != nil {
			return &url.URL{}
		}
		return p
	}
	return &url.URL{}
}

func (c *Config) DefaultApiURL(defaultPort string) string {
	return fmt.Sprintf(DefaultApiURLFormat, c.localIP, defaultPort)
}

func (c *Config) DefaultKernelApiURL(defaultPort string) string {
	return fmt.Sprintf(DefaultApiURLFormat, c.localIP, defaultPort)
}

func (c *Config) RegisterApiURL(id string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.Urls[id]; ok {
		return
	}
	c.Urls[id] = url
}

func (c *Config) RegisterKernelApiURL(id string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.kernelUrls[id]; ok {
		return
	}
	c.kernelUrls[id] = url
}

func (c *Config) LocalIP() string {
	return c.localIP
}

func (c *Config) Debug() bool {
	return c.mode == DebugMode
}
