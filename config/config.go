package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vlog"
)

const urlFormat = "http://%s:%s"

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
}

func New() Config {
	localIP, err := net.LocalIP()
	if err != nil {
		log.Error(err)
		localIP = "127.0.0.1"
	}

	c := Config{
		mode:    ProductionMode,
		localIP: localIP,
		urls: map[string]string{
			"vertex": fmt.Sprintf(urlFormat, localIP, "6130"),
		},
		kernelUrls: map[string]string{
			"vertex": fmt.Sprintf(urlFormat, localIP, "6131"),
		},
	}

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

func (c Config) KernelURL(id string) string {
	if url, ok := c.kernelUrls[id]; ok {
		return url
	}
	log.Error(fmt.Errorf("no url configured for this kernel app"), vlog.String("app_id", id))
	return ""
}

func (c Config) URL(id string) string {
	if url, ok := c.urls[id]; ok {
		return url
	}
	log.Error(fmt.Errorf("no url configured for this app"), vlog.String("app_id", id))
	return ""
}

func (c Config) RegisterApiURL(id string, url string) {
	c.urls[id] = url
}

func (c Config) RegisterKernelApiURL(id string, url string) {
	c.kernelUrls[id] = url
}

func (c Config) LocalIP() string {
	return c.localIP
}

func (c Config) Debug() bool {
	return c.mode == DebugMode
}

func (c Config) Apply() error {
	cfg := ""
	// Only for the non-kernel apps
	for name, url := range c.urls {
		cfg += fmt.Sprintf("window.api_url_%s = \"%s\";\n", name, url)
	}
	return os.WriteFile(path.Join(storage.FSPath, "client", "dist", "config.js"), []byte(cfg), os.ModePerm)
}
