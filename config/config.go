package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
)

const urlFormat = "http://%s:%s"

var Current = New()

type Mode string

const (
	ProductionMode Mode = "production"
	DebugMode      Mode = "debug"
)

type Config struct {
	mode         Mode
	Host         string `json:"host"`
	Ports        map[string]string
	MasterApiKey string `json:"master_api_key"`
}

func New() Config {
	host, err := net.LocalIP()
	if err != nil {
		log.Error(err)
		host = "127.0.0.1"
	}

	// Generate a random master key token
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		log.Error(fmt.Errorf("failed to generate master key token: %w", err))
	}

	c := Config{
		mode: ProductionMode,
		Host: host,
		Ports: map[string]string{
			"VERTEX":        "6130",
			"VERTEX_KERNEL": "6131",
			"VERTEX_PROXY":  "80",
		},
		MasterApiKey: base64.StdEncoding.EncodeToString(token),
	}

	env := os.Environ()
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			key, value := pair[0], pair[1]
			if strings.HasPrefix(key, "VERTEX_PORT_") {
				name := strings.TrimPrefix(key, "VERTEX_PORT_")
				c.Ports[name] = value
			}
		}
	}

	if os.Getenv("DEBUG") == "1" {
		log.Warn("debug mode enabled. proceed with caution!")
		c.mode = DebugMode
	}

	return c
}

func (c Config) VertexURL() string {
	return fmt.Sprintf(urlFormat, c.Host, c.Ports["VERTEX"])
}

func (c Config) KernelURL() string {
	return fmt.Sprintf(urlFormat, c.Host, c.Ports["VERTEX_KERNEL"])
}

func (c Config) GetPort(name string, fallback string) string {
	if port, ok := c.Ports[name]; ok {
		return port
	}
	return fallback
}

func (c Config) Debug() bool {
	return c.mode == DebugMode
}

func (c Config) Apply() error {
	configJsContent := fmt.Sprintf("window.apiURL = \"%s\";", c.Host)
	configJsContent += fmt.Sprintf("window.apiPort_VERTEX = \"%s\";", c.Ports["VERTEX"])
	configJsContent += fmt.Sprintf("window.apiPort_VERTEX_PROXY = \"%s\";", c.Ports["VERTEX_PROXY"])

	return os.WriteFile(path.Join(storage.FSPath, "client", "dist", "config.js"), []byte(configJsContent), os.ModePerm)
}
