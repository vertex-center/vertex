package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/vertex-center/vertex/cmd/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vlog"
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
	urls       map[string]string
	kernelUrls map[string]string
	mu         sync.RWMutex
}

func New() *Config {
	localIP, err := net.LocalIP()
	if err != nil {
		log.Error(err)
		localIP = "127.0.0.1"
	}

	c := &Config{
		mode:    ProductionMode,
		localIP: localIP,
		urls: map[string]string{
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
					c.urls[name] = value
				}
			}
		}
	}

	if os.Getenv("DEBUG") == "1" {
		log.Warn("debug mode enabled. proceed with caution!")
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
			log.Error(err)
			return &url.URL{}
		}
		return p
	}
	log.Error(fmt.Errorf("no url configured for this kernel app"), vlog.String("app_id", id))
	return &url.URL{}
}

func (c *Config) URL(id string) *url.URL {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if u, ok := c.urls[id]; ok {
		p, err := url.Parse(u)
		if err != nil {
			log.Error(err)
			return &url.URL{}
		}
		return p
	}
	log.Error(fmt.Errorf("no url configured for this app"), vlog.String("app_id", id))
	return &url.URL{}
}

func (c *Config) RegisterApiURL(id string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.urls[id]; ok {
		return
	}
	c.urls[id] = url
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

func (c *Config) Apply() error {
	cfg := "window.api_urls = {\n"
	// Only for the non-kernel apps
	for name, u := range c.urls {
		name = strings.ReplaceAll(name, "-", "_")
		cfg += fmt.Sprintf("\t%s: '%s',\n", name, u)
	}
	cfg += "};\n"
	return os.WriteFile(path.Join(storage.FSPath, "client", "dist", "config.js"), []byte(cfg), os.ModePerm)
}
